<template>
  <div class="page-container">
    <p v-if="loadError" class="api-error" role="alert">{{ loadError }}</p>
    <EmptyState v-if="loading" message="Loading departments..." />
    <div class="page-title-section">
      <h2>Departments</h2>
      <p class="breadcrumb">Dashboard > <span>Departments</span></p>
    </div>

    <div class="stats-top-grid">
      <div class="metric-card-box">
        <div class="icon-circle purple-bg">🏢</div>
        <div class="metric-info-data">
          <h4>{{ hasDepartmentFilters ? 'Matching Departments' : 'Total Departments' }}</h4>
          <p class="metric-num-val">{{ filteredDepartments.length }}</p>
          <span class="trend-indicator text-green">{{ hasDepartmentFilters ? 'Matching current filters' : 'Across the organization' }}</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle green-bg">💼</div>
        <div class="metric-info-data">
          <h4>Active Departments</h4>
          <p class="metric-num-val">{{ activeDepartmentCount }}</p>
          <span class="trend-indicator text-blue">In current results</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle blue-bg">👥</div>
        <div class="metric-info-data">
          <h4>Total Employees</h4>
          <p class="metric-num-val">{{ summary.total_employees }}</p>
          <span class="trend-indicator text-gray">Across all departments</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle light-green-bg">👤</div>
        <div class="metric-info-data">
          <h4>Department Heads</h4>
          <p class="metric-num-val">{{ departmentHeadCount }}</p>
          <span class="trend-indicator text-green">In current results</span>
        </div>
      </div>
    </div>

    <div class="table-main-wrapper">
      <div class="controls-action-bar">
        <div class="search-field-container">
          <span class="search-lens">🔍</span>
          <input v-model="searchQuery" type="search" placeholder="Search departments..." class="search-input-element" />
        </div>
        <select v-model="selectedStatus" class="department-status-filter" aria-label="Filter by department status">
          <option value="">All statuses</option>
          <option value="Active">Active</option>
          <option value="Inactive">Inactive</option>
        </select>
        <button class="add-department-btn" @click="startNewDepartment">+ Add Department</button>
      </div>

      <form v-if="showCreateDepartmentForm" class="create-department-form" @submit.prevent="submitNewDepartment">
        <p v-if="formError" class="api-error" role="alert">{{ formError }}</p>
        <h3 class="department-form-title">{{ editingDepartmentId ? 'Edit Department' : 'New Department' }}</h3>
        <div class="form-grid">
          <label>
            Department Name
            <input v-model="newDepartment.name" type="text" placeholder="Enter department name" required />
          </label>
          <label>
            Department Head
            <input v-model="newDepartment.head" type="text" placeholder="Enter department head" />
          </label>
          <label>
            Employees
            <input v-model="newDepartment.employees" type="number" min="0" />
          </label>
          <label>
            Status
            <select v-model="newDepartment.status">
              <option>Active</option>
              <option>Inactive</option>
            </select>
          </label>
        </div>
        <div class="department-form-actions">
          <button type="submit" class="add-department-btn">{{ editingDepartmentId ? 'Update Department' : 'Save Department' }}</button>
          <button type="button" class="cancel-btn" @click="cancelDepartmentForm">Cancel</button>
        </div>
      </form>

      <div class="table-scroll-area">
        <table class="corporate-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Department Name</th>
              <th>Department Head</th>
              <th>Employees</th>
              <th>Status</th>
              <th>Created Date</th>
              <th style="text-align: center;">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="dept in pagedDepartments" :key="dept.id">
              <td class="id-code">{{ dept.id }}</td>
              <td class="dept-name-td">{{ dept.name }}</td>
              <td class="head-name-td">{{ dept.head }}</td>
              <td class="emp-count-td">{{ dept.employees }}</td>
              <td><span class="status-badge-pill" :class="dept.status === 'Active' ? 'active-bg' : 'inactive-bg'">{{ dept.status }}</span></td>
              <td class="date-td">{{ dept.created_date }}</td>
              <td>
                <div class="action-btn-group">
                  <button class="table-icon-btn edit-color" title="Edit department" @click="editDepartment(dept)">📝</button>
                  <button class="table-icon-btn delete-color" title="Delete department" @click="deleteDepartment(dept)">🗑️</button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && !pagedDepartments.length">
              <td colspan="7"><EmptyState :message="filteredDepartments.length ? 'No departments on this page' : 'No departments yet'" /></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-pagination-bar">
        <span class="pagination-text-info">Showing {{ paginationStart }}–{{ paginationEnd }} of {{ filteredDepartments.length }} departments</span>
        <div class="pagination-nav-cluster">
          <button class="nav-arrow-btn" :disabled="currentPage === 1" @click="currentPage--">‹</button>
          <button v-for="page in pageCount" :key="page" class="nav-number-btn" :class="{ 'active-page': page === currentPage }" @click="currentPage = page">{{ page }}</button>
          <button class="nav-arrow-btn" :disabled="currentPage >= pageCount" @click="currentPage++">›</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { apiDelete, apiGet, apiPost, apiPut } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const departments = ref([]);
const summary = ref({ active_departments: 0, total_employees: 0, department_heads: 0 });
const showCreateDepartmentForm = ref(false);
const newDepartment = ref({ name: '', head: '', employees: 0, status: 'Active' });
const searchQuery = ref('');
const selectedStatus = ref('');
const editingDepartmentId = ref(null);
const loadError = ref('');
const loading = ref(false);
const formError = ref('');
const currentPage = ref(1);
const pageSize = 8;
const hasDepartmentFilters = computed(() => Boolean(searchQuery.value.trim() || selectedStatus.value));
const filteredDepartments = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  return departments.value.filter(department => {
    const matchesQuery = !query || [department.name, department.head, department.status]
      .some(value => String(value || '').toLowerCase().includes(query));
    const matchesStatus = !selectedStatus.value || department.status === selectedStatus.value;
    return matchesQuery && matchesStatus;
  });
});
const activeDepartmentCount = computed(() => filteredDepartments.value.filter(department => department.status === 'Active').length);
const departmentHeadCount = computed(() => filteredDepartments.value.filter(department => Boolean(department.head?.trim())).length);
const pageCount = computed(() => Math.max(1, Math.ceil(filteredDepartments.value.length / pageSize)));
const pagedDepartments = computed(() => filteredDepartments.value.slice((currentPage.value - 1) * pageSize, currentPage.value * pageSize));
const paginationStart = computed(() => filteredDepartments.value.length ? (currentPage.value - 1) * pageSize + 1 : 0);
const paginationEnd = computed(() => Math.min(currentPage.value * pageSize, filteredDepartments.value.length));
watch([searchQuery, selectedStatus], () => { currentPage.value = 1; });
watch(filteredDepartments, () => {
  if (currentPage.value > pageCount.value) currentPage.value = pageCount.value;
});

const loadDepartments = async () => {
  loading.value = true;
  loadError.value = '';
  try {
    const data = await apiGet('/api/departments');
    summary.value = {
      active_departments: Number(data.summary?.active_departments) || 0,
      total_employees: Number(data.summary?.total_employees) || 0,
      department_heads: Number(data.summary?.department_heads) || 0,
    };
    departments.value = Array.isArray(data.departments) ? data.departments : [];
    loadError.value = '';
  } catch (error) {
    console.error('Failed to load departments', error);
    loadError.value = error.message;
    departments.value = [];
    summary.value = { active_departments: 0, total_employees: 0, department_heads: 0 };
  } finally {
    loading.value = false;
  }
};

const resetNewDepartment = () => {
  newDepartment.value = { name: '', head: '', employees: 0, status: 'Active' };
};

const cancelDepartmentForm = () => {
  resetNewDepartment();
  editingDepartmentId.value = null;
  formError.value = '';
  showCreateDepartmentForm.value = false;
};

const startNewDepartment = () => {
  cancelDepartmentForm();
  showCreateDepartmentForm.value = true;
};

const editDepartment = (department) => {
  formError.value = '';
  editingDepartmentId.value = department.id;
  newDepartment.value = {
    name: department.name || '',
    head: department.head || '',
    employees: Number(department.employees) || 0,
    status: department.status || 'Active',
  };
  showCreateDepartmentForm.value = true;
};

const submitNewDepartment = async () => {
  formError.value = '';
  try {
    const payload = {
      name: newDepartment.value.name,
      head: newDepartment.value.head,
      employees: Number(newDepartment.value.employees),
      status: newDepartment.value.status,
    };

    if (editingDepartmentId.value !== null) {
      const updated = await apiPut(`/api/departments/${encodeURIComponent(editingDepartmentId.value)}`, payload);
      const index = departments.value.findIndex(department => department.id === editingDepartmentId.value);
      if (index !== -1) departments.value[index] = updated;
    } else {
      const created = await apiPost('/api/departments', payload);
      departments.value.unshift(created);
    }

    resetNewDepartment();
    editingDepartmentId.value = null;
    showCreateDepartmentForm.value = false;
  } catch (error) {
    console.error('Failed to create department', error);
    formError.value = error.message || 'Unable to save this department. Please try again.';
  }
};

const deleteDepartment = async (department) => {
  if (!window.confirm(`Delete ${department.name}?`)) return;
  try {
    await apiDelete(`/api/departments/${encodeURIComponent(department.id)}`);
    departments.value = departments.value.filter(item => item.id !== department.id);
  } catch (error) {
    console.error('Failed to delete department', error);
    alert(error.message);
  }
};

onMounted(loadDepartments);
</script>

<style scoped>
.page-container { background: transparent; animation: fadeInView 0.4s ease-out; }
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; margin: 0 0 16px; }

/* BREADCRUMB HEADER */
.page-title-section { margin-bottom: 20px; }
.page-title-section h2 { margin: 0 0 4px; font-size: 24px; color: #0f172a; font-weight: 700; }
.breadcrumb { font-size: 13px; color: #64748b; margin: 0; }
.breadcrumb span { color: #047857; font-weight: 500; }

/* SUMMARY CARD METRICS GRID */
.stats-top-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 20px; margin-bottom: 25px; }
.metric-card-box { background: white; padding: 20px; border-radius: 12px; display: flex; align-items: center; gap: 15px; border: 1px solid #e2e8f0; box-shadow: 0 1px 3px rgba(0,0,0,0.02); }
.icon-circle { width: 46px; height: 46px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.purple-bg, .green-bg { background-color: #d1fae5; color: #065f46; }
.blue-bg { background-color: #eff6ff; }
.light-green-bg { background-color: #f0fdf4; }
.metric-info-data h4 { margin: 0 0 4px; font-size: 13px; color: #64748b; font-weight: 500; }
.metric-num-val { font-size: 24px; font-weight: 700; margin: 0 0 2px; color: #0f172a; }
.trend-indicator { font-size: 12px; }
.text-green { color: #16a34a; font-weight: 500; }
.text-blue { color: #047857; font-weight: 500; }
.text-gray { color: #94a3b8; }

/* TABLE CARD HOUSING */
.table-main-wrapper { min-width: 0; max-width: 100%; background: white; border-radius: 12px; border: 1px solid #e2e8f0; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.03); padding: 20px; }
.controls-action-bar { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 15px; margin-bottom: 20px; }

/* SEARCH INPUT STRUCTURE */
.search-field-container { position: relative; width: 300px; min-width: 0; max-width: 100%; }
.search-lens { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); color: #94a3b8; font-size: 14px; }
.search-input-element { width: 100%; padding: 9px 14px 9px 38px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: 14px; outline: none; background-color: #f8fafc; transition: 0.2s; }
.search-input-element:focus { border-color: #10b981; background-color: white; box-shadow: 0 0 0 3px rgba(16,185,129,0.12); }
.department-status-filter { min-width: 0; max-width: 100%; min-height: 38px; padding: 8px 12px; border: 1px solid #d7e5df; border-radius: 8px; background: #fff; color: #18352b; font-size: 14px; }
.department-status-filter:focus { border-color: #10b981; outline: 2px solid rgba(16,185,129,0.12); }
.add-department-btn { background-color: #047857; color: white; border: none; padding: 10px 18px; border-radius: 8px; cursor: pointer; font-weight: 600; font-size: 14px; transition: 0.2s; }
.add-department-btn:hover { background-color: #059669; }

.create-department-form {
  margin: 0 0 20px;
  padding: 20px;
  background: #f0f7f4;
  border: 1px solid #d7e5df;
  border-radius: 12px;
}
.department-form-title {
  margin: 0 0 16px;
  color: #18352b;
  font-size: 16px;
  font-weight: 600;
}
.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(min(100%, 210px), 1fr));
  gap: 16px;
  margin-bottom: 20px;
}
.form-grid label {
  display: flex;
  flex-direction: column;
  gap: 7px;
  color: #38594d;
  font-size: 13px;
  font-weight: 600;
}
.form-grid input,
.form-grid select {
  width: 100%;
  min-height: 40px;
  padding: 9px 12px;
  border: 1px solid #cbdcd3;
  border-radius: 8px;
  background: #ffffff;
  color: #18352b;
  font: inherit;
}
.form-grid input:focus,
.form-grid select:focus {
  border-color: #10b981;
  outline: 2px solid rgba(16, 185, 129, 0.14);
}
.department-form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* DATATABLE STYLING */
.table-scroll-area { overflow-x: auto; max-width: 100%; min-width: 0; }
.corporate-table { width: 100%; border-collapse: collapse; font-size: 14px; text-align: left; }
.corporate-table th { background-color: #f8fafc; padding: 14px 16px; font-weight: 600; color: #475569; border-bottom: 2px solid #e2e8f0; font-size: 13px; }
.corporate-table td { padding: 14px 16px; border-bottom: 1px solid #f1f5f9; color: #334155; vertical-align: middle; }

/* DATA ROW ALIGNMENTS */
.id-code { font-family: monospace; font-weight: 600; color: #64748b; }
.dept-name-td { font-weight: 600; color: #0f172a; }
.head-name-td { color: #334155; }
.emp-count-td { font-weight: 500; color: #475569; }
.date-td { color: #64748b; }

@media (max-width: 767px) {
  .controls-action-bar { align-items: stretch; }
  .search-field-container { flex: 1 1 100%; width: 100%; }
  .department-status-filter { flex: 1 1 140px; }
  .add-department-btn { flex: 1 1 100%; }
}

/* STATUS LABELS */
.status-badge-pill { padding: 3px 8px; border-radius: 6px; font-size: 11px; font-weight: 600; text-transform: uppercase; display: inline-block; }
.status-badge-pill.active-bg { background-color: #dcfce7; color: #16a34a; }
.status-badge-pill.inactive-bg { background-color: #f1f5f9; color: #64748b; }

/* TABLE ACTIONS */
.action-btn-group { display: flex; gap: 8px; justify-content: center; }
.table-icon-btn { background: none; border: 1px solid #e2e8f0; padding: 5px 8px; border-radius: 6px; cursor: pointer; transition: 0.2s; font-size: 13px; }
.table-icon-btn:hover { background: #f8fafc; transform: scale(1.05); }

/* PAGINATION STYLES */
.table-pagination-bar { display: flex; justify-content: space-between; align-items: center; margin-top: 20px; flex-wrap: wrap; gap: 15px; border-top: 1px solid #f1f5f9; padding-top: 15px; }
.pagination-text-info { font-size: 13px; color: #64748b; }
.pagination-nav-cluster { display: flex; gap: 5px; }
.nav-arrow-btn, .nav-number-btn { background: white; border: 1px solid #e2e8f0; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; cursor: pointer; font-size: 13px; font-weight: 500; color: #475569; transition: 0.2s; }
.nav-number-btn.active-page { background-color: #047857; color: white; border-color: #047857; }
.nav-arrow-btn:disabled { color: #cbd5e1; cursor: not-allowed; background: #f8fafc; }
.nav-arrow-btn:not(:disabled):hover, .nav-number-btn:not(.active-page):hover { background: #f1f5f9; }

@keyframes fadeInView {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>