<template>
  <div class="page-container">
    <p v-if="loadError" class="api-error" role="alert">{{ loadError }}</p>
    <div class="employee-stats-grid">
      <div class="emp-stat-card border-blue">
        <div class="stat-icon-wrap blue-bg">👥</div>
        <div class="stat-info">
          <h4>All Employees</h4>
          <p class="stat-num">{{ totalEmployees }}</p>
        </div>
      </div>
      <div class="emp-stat-card border-green">
        <div class="stat-icon-wrap green-bg">✅</div>
        <div class="stat-info">
          <h4>Active Staff</h4>
          <p class="stat-num">{{ activeStaff }}</p>
        </div>
      </div>
      <div class="emp-stat-card border-orange">
        <div class="stat-icon-wrap orange-bg">⏳</div>
        <div class="stat-info">
          <h4>On Leave</h4>
          <p class="stat-num">{{ onLeave }}</p>
        </div>
      </div>
    </div>

    <div class="table-controls-card">
      <div class="controls-header">
        <div class="search-box-wrap">
          <span class="search-icon">🔍</span>
          <input v-model="searchQuery" type="search" placeholder="Search employees by name, role or ID..." class="search-input" />
        </div>
        <div class="filter-actions">
          <select v-model="selectedDepartment" class="filter-dropdown">
            <option value="">All Departments</option>
            <option v-for="department in departmentOptions" :key="department" :value="department">{{ department }}</option>
          </select>
          <button class="add-emp-btn" @click="startNewEmployee">+ Add Employee</button>
        </div>
      </div>

      <form v-if="showCreateForm" class="create-employee-form" @submit.prevent="submitNewEmployee">
        <p v-if="formError" class="api-error" role="alert">{{ formError }}</p>
        <div class="form-grid">
          <label>
            Employee ID
            <input v-model="newEmployee.employee_id" type="text" readonly />
          </label>
          <label>
            Full Name
            <input v-model="newEmployee.full_name" type="text" placeholder="Jane Doe" required />
          </label>
          <label>
            Email
            <input v-model="newEmployee.email" type="email" placeholder="jane.doe@example.com" required />
          </label>
          
          <label>
            Password
            <input v-model="newEmployee.password" type="password" :placeholder="editingEmployeeId ? 'Leave blank to keep current password' : 'Enter Login Password'" :required="!editingEmployeeId" />
          </label>
          
          <label>
            Department
            <input v-model="newEmployee.department" type="text" placeholder="Finance" required />
          </label>
          <label>
            Position
            <input v-model="newEmployee.position" type="text" placeholder="Accountant" required />
          </label>
          <label>
            Status
            <select v-model="newEmployee.status">
              <option>Active</option>
              <option>On Leave</option>
              <option>Inactive</option>
            </select>
          </label>
        </div>
        <div class="form-actions">
          <button type="submit" class="submit-btn">{{ editingEmployeeId ? 'Update Employee' : 'Save Employee' }}</button>
          <button type="button" class="cancel-btn" @click="cancelEmployeeForm">Cancel</button>
        </div>
      </form>

      <div class="table-responsive">
        <table class="premium-table">
          <thead>
            <tr>
              <th>Employee ID</th>
              <th>Full Name</th>
              <th>Department</th>
              <th>Position</th>
              <th>Status</th>
              <th style="text-align: center;">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="employee in pagedEmployees" :key="employee.employee_id">
              <td class="emp-id">#{{ employee.employee_id }}</td>
              <td>
                <div class="user-cell">
                  <div class="user-avatar-small">{{ employee.full_name.slice(0, 2).toUpperCase() }}</div>
                  <div>
                    <span class="user-fullname">{{ employee.full_name }}</span>
                    <span class="user-email">{{ employee.email }}</span>
                  </div>
                </div>
              </td>
              <td><span class="dept-tag" :class="departmentTagClass(employee.department)">{{ employee.department }}</span></td>
              <td class="role-text">{{ employee.position }}</td>
              <td><span class="status-pill" :class="employee.status === 'Active' ? 'active' : employee.status === 'On Leave' ? 'on-leave' : 'inactive'">{{ employee.status }}</span></td>
              <td>
                <div class="action-icons">
                  <button class="icon-btn edit-btn" title="Edit employee" @click="editEmployee(employee)">📝</button>
                  <button class="icon-btn delete-btn" title="Delete employee" @click="deleteEmployee(employee)">🗑️</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-pagination">
        <span class="pagination-info">Showing {{ paginationStart }}–{{ paginationEnd }} of {{ filteredEmployees.length }} employees</span>
        <div class="pagination-buttons">
          <button class="page-nav-btn" :disabled="currentPage === 1" @click="currentPage--">Previous</button>
          <button v-for="page in pageCount" :key="page" class="page-num-btn" :class="{ active: page === currentPage }" @click="currentPage = page">{{ page }}</button>
          <button class="page-nav-btn" :disabled="currentPage >= pageCount" @click="currentPage++">Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { apiDelete, apiGet, apiPost, apiPut } from '@/api';

const employees = ref([]);
const showCreateForm = ref(false);
const searchQuery = ref('');
const selectedDepartment = ref('');
const editingEmployeeId = ref('');
const loadError = ref('');
const formError = ref('');
const currentPage = ref(1);
const pageSize = 8;

// 🟢 වෙනස් කරන ලදි: දත්ත ව්‍යුහය ඇතුළට password එක එකතු කරන ලදි
const newEmployee = ref({
  employee_id: '',
  full_name: '',
  email: '',
  password: '', // 👈 මෙන්න මෙතනට ආවා
  department: '',
  position: '',
  status: 'Active',
});

const totalEmployees = computed(() => employees.value.length);
const activeStaff = computed(() => employees.value.filter(emp => emp.status === 'Active').length);
const onLeave = computed(() => employees.value.filter(emp => emp.status === 'On Leave').length);
const departmentOptions = computed(() => [...new Set(employees.value.map(employee => employee.department).filter(Boolean))].sort());
const filteredEmployees = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  return employees.value.filter(employee => {
    const matchesQuery = !query || [employee.employee_id, employee.full_name, employee.position, employee.email]
      .some(value => (value || '').toLowerCase().includes(query));
    const matchesDepartment = !selectedDepartment.value || employee.department === selectedDepartment.value;
    return matchesQuery && matchesDepartment;
  });
});
const pageCount = computed(() => Math.max(1, Math.ceil(filteredEmployees.value.length / pageSize)));
const pagedEmployees = computed(() => filteredEmployees.value.slice((currentPage.value - 1) * pageSize, currentPage.value * pageSize));
const paginationStart = computed(() => filteredEmployees.value.length ? (currentPage.value - 1) * pageSize + 1 : 0);
const paginationEnd = computed(() => Math.min(currentPage.value * pageSize, filteredEmployees.value.length));
watch([searchQuery, selectedDepartment], () => { currentPage.value = 1; });

const departmentTagClass = (department) => {
  const key = (department || '').toLowerCase();
  if (key.includes('it')) return 'it-tag';
  if (key.includes('hr')) return 'hr-tag';
  if (key.includes('finance')) return 'finance-tag';
  return 'dept-tag';
};

const generateEmployeeId = () => {
  const existingIds = new Set(employees.value.map(employee => (employee.employee_id || '').toUpperCase()));
  let nextNumber = employees.value.length + 1;
  let candidate = '';

  do {
    candidate = `EMP${String(nextNumber).padStart(3, '0')}`;
    nextNumber += 1;
  } while (existingIds.has(candidate));

  return candidate;
};

const loadEmployees = async () => {
  try {
    const data = await apiGet('/api/employees');
    employees.value = Array.isArray(data) ? data : [];
  } catch (error) {
    console.error('Failed to load employees', error);
    loadError.value = error.message;
  }
};

// 🟢 වෙනස් කරන ලදි: Reset කරද්දී password එකත් හිස් වෙන්න මෙතනට දැම්මා
const resetNewEmployee = () => {
  newEmployee.value = {
    employee_id: '',
    full_name: '',
    email: '',
    password: '', // 👈 මෙතනටත් ආවා
    department: '',
    position: '',
    status: 'Active',
  };
};

const cancelEmployeeForm = () => {
  resetNewEmployee();
  editingEmployeeId.value = '';
  formError.value = '';
  showCreateForm.value = false;
};

const startNewEmployee = () => {
  cancelEmployeeForm();
  newEmployee.value.employee_id = generateEmployeeId();
  showCreateForm.value = true;
};

const editEmployee = (employee) => {
  formError.value = '';
  editingEmployeeId.value = employee.employee_id;
  newEmployee.value = {
    employee_id: employee.employee_id,
    full_name: employee.full_name || '',
    email: employee.email || '',
    password: '',
    department: employee.department || '',
    position: employee.position || '',
    status: employee.status || 'Active',
  };
  showCreateForm.value = true;
};

const submitNewEmployee = async () => {
  formError.value = '';
  try {
    if (editingEmployeeId.value) {
      const updated = await apiPut(`/api/employees/${encodeURIComponent(editingEmployeeId.value)}`, newEmployee.value);
      const index = employees.value.findIndex(employee => employee.employee_id === editingEmployeeId.value);
      if (index !== -1) employees.value[index] = updated;
    } else {
      const created = await apiPost('/api/employees', newEmployee.value);
      employees.value.unshift(created);
    }
    resetNewEmployee();
    editingEmployeeId.value = '';
    showCreateForm.value = false;
  } catch (error) {
    console.error('Failed to add employee', error);
    if (/already exists|duplicate|unique/i.test(error.message)) {
      formError.value = 'An employee with this ID or email already exists. Please review the values and try again.';
      if (!editingEmployeeId.value) {
        await loadEmployees();
        newEmployee.value.employee_id = generateEmployeeId();
      }
      return;
    }
    formError.value = error.message || 'Unable to save this employee. Please try again.';
  }
};

const deleteEmployee = async (employee) => {
  if (!window.confirm(`Delete ${employee.full_name}?`)) return;
  try {
    await apiDelete(`/api/employees/${encodeURIComponent(employee.employee_id)}`);
    employees.value = employees.value.filter(item => item.employee_id !== employee.employee_id);
  } catch (error) {
    console.error('Failed to delete employee', error);
    alert(error.message);
  }
};

onMounted(loadEmployees);
</script>

<style scoped>
.page-container { background: transparent; animation: fadeIn 0.5s ease; }
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; margin: 0 0 16px; }

/* STATS CARDS GRID */
.employee-stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 20px; margin-bottom: 25px; }
.emp-stat-card { background: white; padding: 20px; border-radius: 14px; display: flex; align-items: center; gap: 15px; border-left: 5px solid #64748b; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.02); }
.border-blue { border-left-color: #047857; }
.border-green { border-left-color: #10b981; }
.border-orange { border-left-color: #10b981; }
.stat-icon-wrap { width: 48px; height: 48px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.blue-bg { background: #eff6ff; }
.green-bg { background: #ecfdf5; }
.orange-bg { background: #fff7ed; }
.stat-info h4 { margin: 0 0 4px; font-size: 13px; color: #64748b; font-weight: 500; text-transform: uppercase; letter-spacing: 0.5px; }
.stat-num { font-size: 26px; font-weight: 700; margin: 0; color: #0f172a; }

/* CONTROLS CARD */
.table-controls-card { background: white; border-radius: 16px; border: 1px solid #e2e8f0; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.04); padding: 20px; }
.controls-header { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 15px; margin-bottom: 20px; }

/* SEARCH BOX */
.search-box-wrap { position: relative; width: 320px; }
.search-icon { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); color: #94a3b8; font-size: 14px; }
.search-input { width: 100%; padding: 10px 14px 10px 40px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: 14px; outline: none; transition: 0.2s; background-color: #f8fafc; }
.search-input:focus { border-color: #10b981; background-color: white; box-shadow: 0 0 0 3px rgba(16,185,129,0.12); }

/* FILTERS & ADD BUTTON */
.filter-actions { display: flex; gap: 12px; align-items: center; }
.filter-dropdown { padding: 10px 14px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: 14px; color: #475569; outline: none; background: white; cursor: pointer; }
.add-emp-btn { background-color: #047857; color: white; border: none; padding: 10px 18px; border-radius: 8px; cursor: pointer; font-weight: 600; font-size: 14px; transition: 0.2s; }
.add-emp-btn:hover { background-color: #059669; transform: translateY(-1px); }

.create-employee-form { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 16px; padding: 20px; margin-bottom: 20px; }
.form-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; margin-bottom: 16px; }
.form-grid label { display: flex; flex-direction: column; gap: 8px; font-size: 14px; color: #334155; }
.form-grid input,
.form-grid select { padding: 12px 14px; border: 1px solid #d1d5db; border-radius: 10px; background: white; font-size: 14px; outline: none; }
.form-actions { display: flex; gap: 12px; flex-wrap: wrap; }
.submit-btn, .cancel-btn { border: none; padding: 11px 18px; border-radius: 10px; font-weight: 600; cursor: pointer; }
.submit-btn { background-color: #10b981; color: white; }
.cancel-btn { background-color: white; color: #475569; border: 1px solid #cbd5e1; }
.submit-btn:hover { background-color: #059669; }
.cancel-btn:hover { background-color: #f8fafc; }

/* PREMIUM TABLE SETUP */
.table-responsive { overflow-x: auto; border: 1px solid #e2e8f0; border-radius: 12px; }
.premium-table { width: 100%; border-collapse: collapse; text-align: left; font-size: 14px; background: white; }
.premium-table th { background-color: #f8fafc; padding: 14px 20px; font-weight: 600; color: #475569; border-bottom: 1px solid #e2e8f0; text-transform: uppercase; font-size: 12px; letter-spacing: 0.5px; }
.premium-table td { padding: 16px 20px; border-bottom: 1px solid #f1f5f9; color: #334155; vertical-align: middle; }

/* TABLE CELL DETAILS */
.emp-id { font-family: 'Courier New', Courier, monospace; font-weight: 700; color: #64748b; }
.user-cell { display: flex; align-items: center; gap: 12px; }
.user-avatar-small { width: 36px; height: 36px; border-radius: 50%; background: #e2e8f0; color: #334155; font-weight: 700; display: flex; align-items: center; justify-content: center; font-size: 13px; }
.user-fullname { display: block; font-weight: 600; color: #0f172a; }
.user-email { display: block; font-size: 12px; color: #64748b; }

/* DEPARTMENT TAGS */
.dept-tag { padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 500; }
.it-tag { background: #d1fae5; color: #065f46; }
.hr-tag { background: #ecfdf5; color: #059669; }
.finance-tag { background: #d1fae5; color: #065f46; }
.role-text { color: #475569; font-weight: 500; }

/* STATUS PILLS */
.status-pill { padding: 4px 10px; border-radius: 100px; font-size: 11px; font-weight: 600; text-transform: uppercase; display: inline-block; }
.status-pill.active { background-color: #dcfce7; color: #16a34a; }

/* ACTIONS STYLE */
.action-icons { display: flex; gap: 8px; justify-content: center; }
.icon-btn { background: none; border: 1px solid #e2e8f0; padding: 6px; border-radius: 6px; cursor: pointer; transition: 0.2s; font-size: 14px; }
.icon-btn:hover { background: #f1f5f9; transform: scale(1.05); }

/* PAGINATION */
.table-pagination { display: flex; justify-content: space-between; align-items: center; margin-top: 20px; flex-wrap: wrap; gap: 15px; }
.pagination-info { font-size: 13px; color: #64748b; }
.pagination-buttons { display: flex; gap: 5px; }
.page-nav-btn, .page-num-btn { background: white; border: 1px solid #e2e8f0; padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 13px; font-weight: 500; color: #475569; transition: 0.2s; }
.page-num-btn.active { background-color: #047857; color: white; border-color: #047857; }
.page-nav-btn:disabled { color: #cbd5e1; cursor: not-allowed; background: #f8fafc; }
.page-nav-btn:not(:disabled):hover, .page-num-btn:not(.active):hover { background: #f1f5f9; }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
