<template>
  <div class="payroll-container">
    <p v-if="loadError" role="alert" class="api-error">{{ loadError }}</p>
    <div class="page-header">
      <div class="header-title">
        <h1>💵 Payroll Management</h1>
        <p class="breadcrumb">Dashboard > Payroll</p>
      </div>
      <div class="header-actions">
        <button class="btn-primary" @click="runPayroll">+ Run Payroll</button>
      </div>
    </div>

    <div class="payroll-cards">
      <div class="p-card">
        <div class="card-icon">💰</div>
        <div class="card-info">
          <h3>Total Payroll (Month)</h3>
          <p class="card-value">{{ summary.total_payroll }}</p>
          <span class="card-sub positive">↑ 4.2% vs last month</span>
        </div>
      </div>

      <div class="p-card">
        <div class="card-icon">✅</div>
        <div class="card-info">
          <h3>Processed</h3>
          <p class="card-value">{{ summary.processed }}</p>
          <span class="card-sub status-done">Fully Paid</span>
        </div>
      </div>

      <div class="p-card">
        <div class="card-icon">⏳</div>
        <div class="card-info">
          <h3>Pending Approvals</h3>
          <p class="card-value">{{ summary.pending_approvals }}</p>
          <span class="card-sub status-pending">Requires Action</span>
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
            <th>ACTION</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in filteredPayrollRecords" :key="record.employee_id">
            <td>{{ record.employee_id }}</td>
            <td>
              <div class="emp-profile">
                <div class="emp-avatar">{{ record.full_name.slice(0, 2).toUpperCase() }}</div>
                <div>
                  <div class="emp-name">{{ record.full_name }}</div>
                  <div class="emp-email">{{ record.email }}</div>
                </div>
              </div>
            </td>
            <td><span class="dept-badge">{{ record.department }}</span></td>
            <td>{{ record.basic_salary }}</td>
            <td>{{ record.allowances }}</td>
            <td>{{ record.deductions }}</td>
            <td class="net-salary">{{ record.net_salary }}</td>
            <td><span class="status-badge" :class="record.status.toLowerCase() === 'paid' ? 'paid' : 'pending'">{{ record.status }}</span></td>
            <td><button class="btn-payslip" @click="downloadPayslip(record)">📄 Payslip</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { apiGet } from '@/api';

const summary = ref({ total_payroll: 'LKR 0', processed: '0 Employees', pending_approvals: '0 Employees' });
const payrollRecords = ref([]);
const loadError = ref('');
const searchQuery = ref('');
const selectedDepartment = ref('');
const departments = computed(() => [...new Set(payrollRecords.value.map(record => record.department).filter(Boolean))].sort());
const filteredPayrollRecords = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  return payrollRecords.value.filter(record => {
    const matchesText = !query || [record.employee_id, record.full_name, record.email].some(value => (value || '').toLowerCase().includes(query));
    return matchesText && (!selectedDepartment.value || record.department === selectedDepartment.value);
  });
});

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

const runPayroll = () => {
  const pending = payrollRecords.value.filter(record => record.status.toLowerCase() !== 'paid');
  if (!pending.length) {
    window.alert('There are no pending payroll records.');
    return;
  }
  if (!window.confirm(`Mark ${pending.length} pending payroll record(s) as paid? This is a demo action and is not saved by the backend.`)) return;
  payrollRecords.value = payrollRecords.value.map(record =>
    record.status.toLowerCase() === 'paid' ? record : { ...record, status: 'PAID' }
  );
  summary.value = { ...summary.value, pending_approvals: '0 Employees' };
};

const loadPayroll = async () => {
  try {
    const data = await apiGet('/api/payroll');
    summary.value = data.summary || summary.value;
    payrollRecords.value = data.payrolls || [];
  } catch (error) {
    console.error('Failed to load payroll', error);
    loadError.value = error.message;
  }
};

onMounted(loadPayroll);
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
  background: white;
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  
  /* 👈 මෙන්න මේ පේළි දෙක අලුතින්ම එකතු කරන්න */
  max-width: 100%;
  overflow-x: auto; /* සමහරවිට screen එක ගොඩක් කුඩා වුණොත් කොටුව ඇතුළේ විතරක් ලස්සනට scroll වෙන්න */
}

.table-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.search-bar {
  padding: 10px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 25px;
  width: 300px;
  outline: none;
  font-size: 13px;
}

.table-filters {
  display: flex;
  gap: 10px;
}

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
</style>