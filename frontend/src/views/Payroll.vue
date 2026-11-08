<template>
  <div class="payroll-container">
    <p v-if="loadError" role="alert" class="api-error">{{ loadError }}</p>
    <p v-if="successMessage" role="status" class="success-message">{{ successMessage }}</p>
    <div class="page-header">
      <div class="header-title">
        <h1>💵 Payroll Management</h1>
        <p class="breadcrumb">Dashboard > Payroll</p>
      </div>
      <button class="btn-primary" type="button" @click="openAddModal">+ Add Payroll</button>
    </div>

    <div class="payroll-cards">
      <div class="p-card">
        <div class="card-icon">💰</div>
        <div class="card-info">
          <h3>Total Payroll (Month)</h3>
          <p class="card-value">{{ formatMoney(summary.total_payroll) }}</p>
        </div>
      </div>

      <div class="p-card">
        <div class="card-icon">✅</div>
        <div class="card-info">
          <h3>Processed</h3>
          <p class="card-value">{{ formatCount(summary.processed) }}</p>
        </div>
      </div>

      <div class="p-card">
        <div class="card-icon">⏳</div>
        <div class="card-info">
          <h3>Pending Approvals</h3>
          <p class="card-value">{{ formatCount(summary.pending_approvals) }}</p>
        </div>
      </div>
    </div>

    <div class="table-section">
      <div class="table-header">
        <input v-model="searchQuery" type="search" placeholder="Search employee payroll..." class="search-bar" />
        <div class="table-filters">
          <select v-model="selectedDepartment" class="filter-select">
            <option value="">All Departments</option>
            <option v-for="department in departments" :key="department" :value="department">{{ department }}</option>
          </select>
          <button class="btn-secondary" @click="exportPayroll">Export Report</button>
        </div>
      </div>

      <p v-if="employeeLoadError" role="alert" class="api-error">{{ employeeLoadError }}</p>
      <div class="payroll-table-scroll">
        <table class="payroll-table">
        <thead>
          <tr>
            <th>EMP ID</th>
            <th>EMPLOYEE</th>
            <th>DEPARTMENT</th>
            <th>BASIC SALARY</th>
            <th>ALLOWANCES</th>
            <th>DEDUCTIONS</th>
            <th>NET SALARY</th>
            <th>STATUS</th>
            <th>ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in filteredPayrollRecords" :key="record.id">
            <td>{{ record.employee_id }}</td>
            <td>
              <div class="emp-profile">
                <div class="emp-avatar">{{ (record.full_name || '').slice(0, 2).toUpperCase() }}</div>
                <div>
                  <div class="emp-name">{{ record.full_name }}</div>
                  <div class="emp-email">{{ record.email }}</div>
                </div>
              </div>
            </td>
            <td><span class="dept-badge">{{ record.department }}</span></td>
            <td>{{ formatMoney(record.basic_salary) }}</td>
            <td>{{ formatMoney(record.allowances) }}</td>
            <td>{{ formatMoney(record.deductions) }}</td>
            <td class="net-salary">{{ formatMoney(record.net_salary) }}</td>
            <td><span class="status-badge" :class="(record.status || '').toLowerCase() === 'paid' ? 'paid' : 'pending'">{{ record.status || 'Pending' }}</span></td>
            <td>
              <div class="payroll-actions">
                <button class="icon-action edit-action" type="button" title="Edit payroll" aria-label="Edit payroll" @click="openEditModal(record)">✏️</button>
                <button v-if="(record.status || '').toLowerCase() === 'pending'" class="action-button pay-action" type="button" :disabled="busyRecordId === record.id" @click="markAsPaid(record)">Mark as Paid</button>
                <button class="icon-action delete-action" type="button" title="Delete payroll" aria-label="Delete payroll" :disabled="busyRecordId === record.id" @click="deletePayroll(record)">🗑️</button>
                <button class="btn-payslip" type="button" @click="downloadPayslip(record)">📄 Payslip</button>
              </div>
            </td>
          </tr>
          <tr v-if="!loading && !filteredPayrollRecords.length">
            <td colspan="9"><EmptyState :message="payrollRecords.length ? 'No payroll records match these filters' : 'No payroll records yet'" /></td>
          </tr>
          <tr v-if="loading">
            <td colspan="9"><EmptyState message="Loading payroll..." /></td>
          </tr>
        </tbody>
        </table>
      </div>
    </div>

    <div v-if="showPayrollModal" class="modal-backdrop" @click.self="closePayrollModal">
      <section class="payroll-modal" role="dialog" aria-modal="true" :aria-labelledby="'payroll-modal-title'">
        <h2 id="payroll-modal-title">{{ editingRecord ? 'Edit Payroll' : 'Add Payroll' }}</h2>
        <p v-if="formError" role="alert" class="api-error">{{ formError }}</p>
        <form @submit.prevent="savePayroll">
          <label class="modal-field">
            Employee
            <select v-model="form.employee_id" :disabled="Boolean(editingRecord)" :required="!editingRecord">
              <option value="" disabled>Select employee</option>
              <option v-if="employeesLoading" value="" disabled>Loading employees...</option>
              <option v-else-if="!employees.length" value="" disabled>No employees available</option>
              <option v-for="employee in employees" :key="employee.employee_id" :value="employee.employee_id">
                {{ employee.employee_id }} - {{ employee.full_name }}
              </option>
            </select>
            <span v-if="validationErrors.employee_id" class="field-error">{{ validationErrors.employee_id }}</span>
          </label>
          <label class="modal-field">
            Basic Salary
            <input v-model.number="form.basic_salary" type="number" min="0" step="any" required />
            <span v-if="validationErrors.basic_salary" class="field-error">{{ validationErrors.basic_salary }}</span>
          </label>
          <label class="modal-field">
            Allowances
            <input v-model.number="form.allowances" type="number" min="0" step="any" required />
            <span v-if="validationErrors.allowances" class="field-error">{{ validationErrors.allowances }}</span>
          </label>
          <label class="modal-field">
            Deductions
            <input v-model.number="form.deductions" type="number" min="0" step="any" required />
            <span v-if="validationErrors.deductions" class="field-error">{{ validationErrors.deductions }}</span>
          </label>
          <label class="modal-field">
            Net Salary
            <input :value="formatMoney(netSalary)" type="text" readonly />
          </label>
          <div class="modal-actions">
            <button class="btn-primary" type="submit" :disabled="saving">{{ saving ? 'Saving...' : editingRecord ? 'Save Changes' : 'Add Payroll' }}</button>
            <button class="modal-cancel" type="button" :disabled="saving" @click="closePayrollModal">Cancel</button>
          </div>
        </form>
      </section>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { apiDelete, apiGet, apiPost, apiPut } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const summary = ref({ total_payroll: 0, processed: 0, pending_approvals: 0 });
const payrollRecords = ref([]);
const loadError = ref('');
const loading = ref(false);
const searchQuery = ref('');
const selectedDepartment = ref('');
const employees = ref([]);
const employeesLoading = ref(false);
const employeeLoadError = ref('');
const showPayrollModal = ref(false);
const editingRecord = ref(null);
const saving = ref(false);
const busyRecordId = ref(null);
const formError = ref('');
const successMessage = ref('');
const validationErrors = ref({});
const form = ref({ employee_id: '', basic_salary: 0, allowances: 0, deductions: 0 });
const departments = computed(() => [...new Set(payrollRecords.value.map(record => record.department).filter(Boolean))].sort());
const filteredPayrollRecords = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  return payrollRecords.value.filter(record => {
    const matchesText = !query || [record.employee_id, record.full_name, record.email].some(value => (value || '').toLowerCase().includes(query));
    return matchesText && (!selectedDepartment.value || record.department === selectedDepartment.value);
  });
  const netSalary = computed(() => toNumber(form.value.basic_salary) + toNumber(form.value.allowances) - toNumber(form.value.deductions));
});

const toNumber = value => {
  const number = Number(value ?? 0);
  return Number.isFinite(number) ? number : 0;
};
const formatMoney = value => `LKR ${toNumber(value).toLocaleString(undefined, { maximumFractionDigits: 2 })}`;
const formatCount = value => toNumber(value).toLocaleString();

const resetForm = () => {
  form.value = { employee_id: '', basic_salary: 0, allowances: 0, deductions: 0 };
  validationErrors.value = {};
  formError.value = '';
  editingRecord.value = null;
};

const openAddModal = () => {
  resetForm();
  successMessage.value = '';
  showPayrollModal.value = true;
};

const openEditModal = record => {
  resetForm();
  editingRecord.value = record;
  form.value = {
    employee_id: record.employee_id || '',
    basic_salary: toNumber(record.basic_salary),
    allowances: toNumber(record.allowances),
    deductions: toNumber(record.deductions),
  };
  successMessage.value = '';
  showPayrollModal.value = true;
};

const closePayrollModal = () => {
  showPayrollModal.value = false;
  resetForm();
};

const validateForm = () => {
  const errors = {};
  if (!editingRecord.value && !form.value.employee_id) errors.employee_id = 'Select an employee.';
  if (toNumber(form.value.basic_salary) <= 0) errors.basic_salary = 'Basic salary must be greater than 0.';
  if (toNumber(form.value.allowances) < 0) errors.allowances = 'Allowances cannot be negative.';
  if (toNumber(form.value.deductions) < 0) errors.deductions = 'Deductions cannot be negative.';
  if (toNumber(form.value.deductions) > toNumber(form.value.basic_salary) + toNumber(form.value.allowances)) {
    errors.deductions = 'Deductions cannot exceed basic salary plus allowances.';
  }
  validationErrors.value = errors;
  return Object.keys(errors).length === 0;
};

const savePayroll = async () => {
  formError.value = '';
  if (!validateForm()) return;
  saving.value = true;
  try {
    const payload = {
      basic_salary: toNumber(form.value.basic_salary),
      allowances: toNumber(form.value.allowances),
      deductions: toNumber(form.value.deductions),
    };
    if (editingRecord.value) {
      await apiPut(`/api/payroll/${editingRecord.value.id}`, payload);
    } else {
      await apiPost('/api/payroll', { employee_id: form.value.employee_id, ...payload });
    }
    const message = editingRecord.value ? 'Payroll updated.' : 'Payroll added.';
    closePayrollModal();
    await loadPayroll();
    successMessage.value = message;
  } catch (error) {
    console.error('Failed to save payroll', error);
    formError.value = error.message || 'Unable to save payroll.';
  } finally {
    saving.value = false;
  }
};

const markAsPaid = async record => {
  if (!window.confirm(`Mark payroll for ${record.full_name} as paid?`)) return;
  busyRecordId.value = record.id;
  loadError.value = '';
  successMessage.value = '';
  try {
    await apiPut(`/api/payroll/${record.id}/pay`, {});
    await loadPayroll();
    successMessage.value = 'Payroll marked as paid.';
  } catch (error) {
    console.error('Failed to mark payroll as paid', error);
    loadError.value = error.message;
  } finally {
    busyRecordId.value = null;
  }
};

const deletePayroll = async record => {
  if (!window.confirm(`Delete payroll for ${record.full_name}?`)) return;
  busyRecordId.value = record.id;
  loadError.value = '';
  successMessage.value = '';
  try {
    await apiDelete(`/api/payroll/${record.id}`);
    await loadPayroll();
    successMessage.value = 'Payroll deleted.';
  } catch (error) {
    console.error('Failed to delete payroll', error);
    loadError.value = error.message;
  } finally {
    busyRecordId.value = null;
  }
};

const downloadFile = (name, contents, type = 'text/csv;charset=utf-8') => {
  const url = URL.createObjectURL(new Blob([contents], { type }));
  const link = document.createElement('a');
  link.href = url;
  link.download = name;
  link.click();
  URL.revokeObjectURL(url);
};

const toCsv = (records) => {
  if (!records.length) return '';
  const columns = Object.keys(records[0]);
  const escape = value => `"${String(value ?? '').replaceAll('"', '""')}"`;
  return [columns.join(','), ...records.map(record => columns.map(column => escape(record[column])).join(','))].join('\r\n');
};

const exportPayroll = () => {
  downloadFile('payroll-report.csv', toCsv(filteredPayrollRecords.value));
};

const downloadPayslip = (record) => {
  const contents = `Payslip for ${record.full_name}\r\n${toCsv([record])}`;
  downloadFile(`payslip-${record.employee_id}.txt`, contents, 'text/plain;charset=utf-8');
};

const loadPayroll = async () => {
  loading.value = true;
  loadError.value = '';
  try {
    const data = await apiGet('/api/payroll');
    summary.value = {
      total_payroll: toNumber(data.summary?.total_payroll),
      processed: toNumber(data.summary?.processed),
      pending_approvals: toNumber(data.summary?.pending_approvals),
    };
    payrollRecords.value = Array.isArray(data.payrolls) ? data.payrolls : [];
  } catch (error) {
    console.error('Failed to load payroll', error);
    loadError.value = error.message;
    summary.value = { total_payroll: 0, processed: 0, pending_approvals: 0 };
    payrollRecords.value = [];
  } finally {
    loading.value = false;
  }
};

const loadEmployees = async () => {
  employeesLoading.value = true;
  employeeLoadError.value = '';
  try {
    const data = await apiGet('/api/employees');
    employees.value = Array.isArray(data) ? data : [];
  } catch (error) {
    console.error('Failed to load employees for payroll', error);
    employeeLoadError.value = `Unable to load employees: ${error.message}`;
    employees.value = [];
  } finally {
    employeesLoading.value = false;
  }
};

onMounted(() => {
  loadPayroll();
  loadEmployees();
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

.payroll-container {
  padding: 25px;
  font-family: 'Inter', sans-serif;
  background-color: #f8fafc;
  min-height: 100vh;
}
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; }
.success-message { color: #047857; background: #ecfdf5; border: 1px solid #a7f3d0; border-radius: 8px; padding: 10px 14px; }

/* Header Styles */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
}

.header-title h1 {
  font-size: 24px;
  font-weight: 700;
  color: #062c22; /* ඔයාගේ තද කොළ පාට */
  margin: 0;
}

.breadcrumb {
  font-size: 12px;
  color: #64748b;
  margin: 4px 0 0;
}

.btn-primary {
  background-color: #10b981; /* මින්ට් කොළ/ලස්සන කොළ පාට */
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 25px;
  font-weight: 500;
  cursor: pointer;
  transition: 0.2s;
}

.btn-primary:hover {
  background-color: #059669;
}
.btn-primary:disabled, .modal-cancel:disabled, .icon-action:disabled, .action-button:disabled { opacity: 0.6; cursor: not-allowed; }

/* Cards Styles */
.payroll-cards {
  display: flex;
  gap: 20px;
  margin-bottom: 25px;
}

.p-card {
  flex: 1;
  background: white;
  border-radius: 15px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.card-icon {
  font-size: 28px;
  background: #f0fdf4;
  padding: 12px;
  border-radius: 12px;
}

.card-info h3 {
  font-size: 13px;
  color: #64748b;
  margin: 0 0 4px;
  font-weight: 500;
}

.card-value {
  font-size: 20px;
  font-weight: 700;
  color: #062c22;
  margin: 0 0 4px;
}

.card-sub {
  font-size: 11px;
}
.positive { color: #10b981; }
.status-done { color: #059669; font-weight: 600; }
.status-pending { color: #047857; font-weight: 600; }

/* Table Section Styles */
.table-section {
  min-width: 0;
  background: white;
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  
  /* 👈 මෙන්න මේ පේළි දෙක අලුතින්ම එකතු කරන්න */
  max-width: 100%;
}

.table-header {
  display: flex;
  min-width: 0;
  max-width: 100%;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.search-bar {
  padding: 10px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 25px;
  flex: 1 1 260px;
  min-width: 0;
  max-width: 100%;
  width: 300px;
  box-sizing: border-box;
  outline: none;
  font-size: 13px;
}

.table-filters {
  display: flex;
  flex-wrap: wrap;
  min-width: 0;
  max-width: 100%;
  gap: 10px;
}

.filter-select { min-width: 0; max-width: 100%; }
.payroll-table-scroll { width: 100%; max-width: 100%; min-width: 0; overflow-x: auto; }

.filter-select {
  padding: 0 15px;
  border: 1px solid #e2e8f0;
  border-radius: 25px;
  outline: none;
  font-size: 13px;
  color: #475569;
}

.btn-secondary {
  background: white;
  border: 1px solid #cbd5e1;
  padding: 10px 20px;
  border-radius: 25px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: #475569;
}

/* Table Design */
/* 🟢 ටේබල් එකේ මූලික සැකසුම (Auto Layout එකට හැරවීම) */
.payroll-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 13px;
  table-layout: auto; /* 👈 'fixed' වෙනුවට 'auto' දැම්මා, දැන් අකුරු එක උඩ හැප්පෙන්නේ නැහැ */
}

/* 🟢 අකුරු වලට ලස්සනට හුස්ම ගන්න ඉඩ දීලා Padding සකස් කිරීම */
.payroll-table th,
.payroll-table td {
  padding: 12px 10px; /* 👈 උඩ-යට සහ දෙපැත්තට පරතරය නියම ගාණට හැදුවා */
  white-space: nowrap; /* 👈 ගණන් සහ අකුරු පේළි කැඩිලා යටට යෑම වළක්වයි */
}

/* 🟢 නම තියෙන කොටසට විතරක් ලස්සනට වැඩිපුර ඉඩක් තබා ගැනීම */
.payroll-table td:nth-child(2) {
  min-width: 180px;
}

.emp-profile {
  display: flex;
  align-items: center;
  gap: 10px;
}

.emp-avatar {
  background: #e2e8f0;
  color: #475569;
  font-weight: 600;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
}

.emp-name {
  font-weight: 600;
  color: #062c22;
}

.emp-email {
  font-size: 11px;
  color: #64748b;
}

.dept-badge {
  background: #eff6ff;
  color: #047857;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.net-salary {
  font-weight: 700;
  color: #062c22;
}

/* Status Badges */
.status-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 10px;
  font-weight: 700;
}

.status-badge.paid {
  background: #dcfce7;
  color: #15803d;
}

.status-badge.pending {
  background: #fef3c7;
  color: #b45309;
}

.btn-payslip {
  background: #f1f5f9;
  border: none;
  color: #475569;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: 0.2s;
}

.btn-payslip:hover {
  background: #e2e8f0;
  color: #062c22;
}

.payroll-actions { display: flex; align-items: center; gap: 6px; }
.icon-action, .action-button, .modal-cancel {
  min-height: 34px;
  padding: 6px 9px;
  border: 1px solid #cbd5e1;
  border-radius: 7px;
  background: white;
  cursor: pointer;
  white-space: nowrap;
}
.edit-action { color: #047857; }
.delete-action { color: #b91c1c; }
.pay-action { color: #047857; font-size: 11px; font-weight: 600; }
.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: rgba(15, 23, 42, 0.48);
}
.payroll-modal {
  width: min(100%, 520px);
  max-height: 90vh;
  overflow-y: auto;
  padding: 26px;
  border-radius: 16px;
  background: white;
  box-shadow: 0 20px 50px rgba(15, 23, 42, 0.2);
}
.payroll-modal h2 { margin: 0 0 20px; color: #062c22; font-size: 20px; }
.payroll-modal form { display: grid; gap: 14px; }
.modal-field { display: grid; gap: 6px; color: #334155; font-size: 13px; font-weight: 600; }
.modal-field input, .modal-field select {
  width: 100%;
  min-height: 40px;
  padding: 9px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: white;
  color: #0f172a;
  font: inherit;
}
.modal-field select:disabled, .modal-field input[readonly] { background: #f1f5f9; }
.field-error { color: #b91c1c; font-size: 12px; font-weight: 400; }
.modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 6px; }
.modal-cancel { padding: 9px 16px; }

@media (max-width: 760px) {
  .page-header, .table-header { align-items: flex-start; flex-direction: column; gap: 12px; }
  .payroll-cards { flex-direction: column; }
  .search-bar, .table-filters { width: 100%; max-width: 100%; }
  .search-bar { flex: 1 1 100%; }
  .table-filters > * { min-width: 0; max-width: 100%; }
}
</style>