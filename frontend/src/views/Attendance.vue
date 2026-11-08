<template>
  <div class="page-container">
    <p v-if="loadError" class="api-error" role="alert">{{ loadError }}</p>
    <EmptyState v-if="loading" message="Loading attendance..." />
    <div class="page-title-section">
      <h2>Attendance</h2>
      <p class="breadcrumb">Dashboard > <span>Attendance</span></p>
    </div>

    <div class="stats-top-grid">
      <div class="metric-card-box">
        <div class="icon-circle green-bg">💚</div>
        <div class="metric-info-data">
          <h4>Present Today</h4>
          <p class="metric-num-val">{{ summary.present_today }}</p>
          <span class="trend-indicator text-green">● {{ attendanceRate }}%</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle red-bg">❤️</div>
        <div class="metric-info-data">
          <h4>Absent Today</h4>
          <p class="metric-num-val">{{ summary.absent_today }}</p>
          <span class="trend-indicator text-red">● {{ absentRate }}%</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle orange-bg">💛</div>
        <div class="metric-info-data">
          <h4>Late Comers</h4>
          <p class="metric-num-val">{{ summary.late_comers }}</p>
          <span class="trend-indicator text-orange">● {{ lateRate }}%</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle blue-bg">💙</div>
        <div class="metric-info-data">
          <h4>Total Employees</h4>
          <p class="metric-num-val">{{ summary.total_employees }}</p>
        </div>
      </div>
    </div>

    <div class="attendance-split-grid">
      
      <div class="table-main-wrapper">
        <div class="controls-action-bar">
          <div class="filters-left-cluster">
            <div class="date-picker-mock">📅 {{ formattedCurrentDate }}</div>
            
            <select class="dropdown-filter" v-model="selectedDepartment">
              <option value="">All Departments</option>
              <option v-for="dept in departmentsList" :key="dept.id" :value="dept.name">
                {{ dept.name }}
              </option>
            </select>

            <select class="dropdown-filter" v-model="selectedStatus">
              <option value="">All Status</option>
              <option value="Present">Present</option>
              <option value="Absent">Absent</option>
              <option value="Late">Late</option>
            </select>
          </div>
          <div class="action-buttons-right">
            <button type="button" class="add-record-btn" @click="startNewAttendance">+ Add Attendance</button>
            <button type="button" class="filter-btn" @click="applyFilters">Apply Filters</button>
            <button v-if="hasFilters" type="button" class="clear-filter-btn" @click="clearFilters">Clear Filters</button>
          </div>
        </div>

        <form v-if="showCreateAttendanceForm" class="create-employee-form" style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin-bottom: 20px;" @submit.prevent="submitNewAttendance">
          <div class="form-grid" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 15px;">
            <label style="display: flex; flex-direction: column; gap: 6px; font-size: 13px; font-weight: 600;">
              Employee Name
              <input v-model="newAttendance.employee" type="text" placeholder="Enter name" required style="padding: 10px; border: 1px solid #cbd5e1; border-radius: 8px;" />
            </label>
            <label style="display: flex; flex-direction: column; gap: 6px; font-size: 13px; font-weight: 600;">
              Department
              <input v-model="newAttendance.department" type="text" placeholder="e.g. IT" required style="padding: 10px; border: 1px solid #cbd5e1; border-radius: 8px;" />
            </label>
            <label style="display: flex; flex-direction: column; gap: 6px; font-size: 13px; font-weight: 600;">
              Check In
              <input v-model="newAttendance.check_in" type="text" placeholder="09:00 AM" style="padding: 10px; border: 1px solid #cbd5e1; border-radius: 8px;" />
            </label>
            <label style="display: flex; flex-direction: column; gap: 6px; font-size: 13px; font-weight: 600;">
              Check Out
              <input v-model="newAttendance.check_out" type="text" placeholder="05:00 PM" style="padding: 10px; border: 1px solid #cbd5e1; border-radius: 8px;" />
            </label>
            <label style="display: flex; flex-direction: column; gap: 6px; font-size: 13px; font-weight: 600;">
              Status
              <select v-model="newAttendance.status" style="padding: 10px; border: 1px solid #cbd5e1; border-radius: 8px; background: white;">
                <option>Present</option>
                <option>Absent</option>
                <option>Late</option>
              </select>
            </label>
          </div>
          <div class="form-actions" style="margin-top: 15px; display: flex; gap: 10px;">
            <button type="submit" class="submit-btn" style="background-color: #10b981; color: white; border: none; padding: 10px 16px; border-radius: 8px; font-weight: 600; cursor: pointer;">{{ editingAttendanceId ? 'Update Attendance' : 'Save Attendance' }}</button>
            <button type="button" class="cancel-btn" @click="cancelAttendanceForm" style="background-color: white; color: #475569; border: 1px solid #cbd5e1; padding: 10px 16px; border-radius: 8px; font-weight: 600; cursor: pointer;">Cancel</button>
          </div>
        </form>

        <div class="table-scroll-area">
          <table class="corporate-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Employee</th>
                <th>Department</th>
                <th>Check In</th>
                <th>Check Out</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in pagedAttendanceRecords" :key="record.id">
                <td class="id-code">#{{ record.id }}</td>
                <td class="emp-name-td">{{ record.employee }}</td>
                <td><span class="dept-label-p" :class="departmentTagClass(record.department)">{{ record.department }}</span></td>
                <td>{{ record.check_in }}</td>
                <td>{{ record.check_out }}</td>
                <td><span class="status-badge-pill" :class="statusClass(record.status)">{{ record.status }}</span></td>
                <td>
                  <button class="table-icon-btn edit-color" title="Edit attendance" @click="editAttendance(record)">📝</button>
                  <button class="table-icon-btn delete-color" title="Delete attendance" @click="deleteAttendance(record)">🗑️</button>
                </td>
              </tr>
              <tr v-if="!pagedAttendanceRecords.length">
                <td colspan="7"><EmptyState :message="attendanceRecords.length ? 'No attendance records match these filters' : 'No attendance records yet'" /></td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="table-pagination-bar">
          <span class="pagination-text-info">Showing {{ paginationStart }}–{{ paginationEnd }} of {{ filteredAttendanceRecords.length }} records</span>
          <div class="pagination-nav-cluster">
            <button class="nav-arrow-btn" :disabled="currentPage === 1" @click="currentPage--">‹</button>
            <button v-for="page in pageCount" :key="page" class="nav-number-btn" :class="{ 'active-page': page === currentPage }" @click="currentPage = page">{{ page }}</button>
            <button class="nav-arrow-btn" :disabled="currentPage >= pageCount" @click="currentPage++">›</button>
          </div>
        </div>
      </div>

      <div class="analytics-side-panel">
        <div class="chart-card">
          <h3>Attendance Overview</h3>
          <div class="pie-container">
            <canvas v-if="attendanceRecords.length" id="attendanceDoughnutChart"></canvas>
            <EmptyState v-else :message="loading ? 'Loading attendance...' : 'No attendance data yet'" height="100%" />
          </div>
        </div>

        <div class="chart-card department-progress-card">
          <h3>Department Wise (Present %)</h3>
          <div class="progress-list">
            <div v-for="department in departmentAttendance" :key="department.name" class="progress-item">
              <div class="progress-info-text"><span>{{ department.name }}</span><span>{{ department.rate }}%</span></div>
              <div class="bar-bg"><div class="bar-fill" :style="{ width: `${department.rate}%`, backgroundColor: '#047857' }"></div></div>
            </div>
            <EmptyState v-if="!departmentAttendance.length" message="No department attendance data yet" />
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, nextTick, ref, watch } from 'vue';
import Chart from 'chart.js/auto';
import { apiDelete, apiGet, apiPost, apiPut } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const summary = ref({ present_today: 0, absent_today: 0, late_comers: 0, total_employees: 0 });
const attendanceRecords = ref([]);
const showCreateAttendanceForm = ref(false);
const newAttendance = ref({ employee: '', department: '', check_in: '', check_out: '', status: 'Present' });
const editingAttendanceId = ref(null);
const loadError = ref('');
const loading = ref(false);

const departmentsList = ref([]);
const selectedDepartment = ref('');

const selectedStatus = ref('');
const activeDepartmentFilter = ref('');
const activeStatusFilter = ref('');
const currentPage = ref(1);
const pageSize = 8;
const hasFilters = computed(() => Boolean(
  selectedDepartment.value || selectedStatus.value || activeDepartmentFilter.value || activeStatusFilter.value
));

let chartInstance = null;

const formattedCurrentDate = computed(() => {
  const today = new Date();
  return today.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
});

const filteredAttendanceRecords = computed(() => {
  return attendanceRecords.value.filter(record => {
    const matchesDepartment = !activeDepartmentFilter.value ||
      (record.department || '').toLowerCase() === activeDepartmentFilter.value.toLowerCase();
      
    const matchesStatus = !activeStatusFilter.value ||
      record.status === activeStatusFilter.value;
      
    return matchesDepartment && matchesStatus;
  });
});
const pageCount = computed(() => Math.max(1, Math.ceil(filteredAttendanceRecords.value.length / pageSize)));
const pagedAttendanceRecords = computed(() => filteredAttendanceRecords.value.slice((currentPage.value - 1) * pageSize, currentPage.value * pageSize));
const paginationStart = computed(() => filteredAttendanceRecords.value.length ? (currentPage.value - 1) * pageSize + 1 : 0);
const paginationEnd = computed(() => Math.min(currentPage.value * pageSize, filteredAttendanceRecords.value.length));
watch([activeDepartmentFilter, activeStatusFilter], () => { currentPage.value = 1; });

const applyFilters = () => {
  activeDepartmentFilter.value = selectedDepartment.value;
  activeStatusFilter.value = selectedStatus.value;
};

const clearFilters = () => {
  selectedDepartment.value = '';
  selectedStatus.value = '';
  activeDepartmentFilter.value = '';
  activeStatusFilter.value = '';
};

const attendanceRate = computed(() => {
  if (!summary.value.total_employees) return 0;
  return ((summary.value.present_today / summary.value.total_employees) * 100).toFixed(1);
});
const absentRate = computed(() => {
  if (!summary.value.total_employees) return 0;
  return ((summary.value.absent_today / summary.value.total_employees) * 100).toFixed(1);
});
const lateRate = computed(() => {
  if (!summary.value.total_employees) return 0;
  return ((summary.value.late_comers / summary.value.total_employees) * 100).toFixed(1);
});
const departmentAttendance = computed(() => {
  const grouped = new Map();
  for (const record of attendanceRecords.value) {
    const name = record.department?.trim();
    if (!name) continue;
    const current = grouped.get(name) || { total: 0, present: 0 };
    current.total += 1;
    if ((record.status || '').toLowerCase() === 'present') current.present += 1;
    grouped.set(name, current);
  }
  return [...grouped].map(([name, counts]) => ({
    name,
    rate: counts.total ? ((counts.present / counts.total) * 100).toFixed(1) : '0.0',
  }));
});

const departmentTagClass = (department) => {
  const key = (department || '').toLowerCase();
  if (key.includes('it') || key.includes('information')) return 'blue-tag';
  if (key.includes('hr') || key.includes('human')) return 'purple-tag';
  if (key.includes('finance')) return 'orange-tag';
  if (key.includes('marketing')) return 'pink-tag';
  return 'blue-tag';
};

const statusClass = (status) => {
  if (status === 'Present') return 'present-bg';
  if (status === 'Late') return 'late-bg';
  if (status === 'Absent') return 'absent-bg';
  return 'present-bg';
};

const attendanceDistribution = computed(() => {
  const counts = new Map();
  for (const record of attendanceRecords.value) {
    const status = record.status || 'Unknown';
    counts.set(status, (counts.get(status) || 0) + 1);
  }
  return [...counts].map(([label, count]) => ({ label, count }));
});

const updateChartData = () => {
  if (chartInstance) {
    chartInstance.data.labels = attendanceDistribution.value.map(item => item.label);
    chartInstance.data.datasets[0].data = attendanceDistribution.value.map(item => item.count);
    chartInstance.update();
  }
};

const renderAttendanceChart = async () => {
  await nextTick();
  if (!attendanceRecords.value.length) {
    chartInstance?.destroy();
    chartInstance = null;
    return;
  }
  const ctx = document.getElementById('attendanceDoughnutChart');
  if (!ctx) return;
  if (chartInstance) {
    updateChartData();
    return;
  }
  chartInstance = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels: attendanceDistribution.value.map(item => item.label),
      datasets: [{
        data: attendanceDistribution.value.map(item => item.count),
        backgroundColor: ['#047857', '#10b981', '#a7f3d0'],
        borderWidth: 2
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { position: 'right', labels: { boxWidth: 12, font: { size: 12 } } }
      }
    }
  });
};

const loadAttendance = async (showLoading = true) => {
  if (showLoading) loading.value = true;
  loadError.value = '';
  try {
    const data = await apiGet('/api/attendance');
    summary.value = {
      present_today: Number(data.summary?.present_today) || 0,
      absent_today: Number(data.summary?.absent_today) || 0,
      late_comers: Number(data.summary?.late_comers) || 0,
      total_employees: Number(data.summary?.total_employees) || 0,
    };
    attendanceRecords.value = Array.isArray(data.attendance) ? data.attendance : [];
    await renderAttendanceChart();
  } catch (error) {
    console.error('Failed to load attendance', error);
    loadError.value = error.message;
    summary.value = { present_today: 0, absent_today: 0, late_comers: 0, total_employees: 0 };
    attendanceRecords.value = [];
    await renderAttendanceChart();
  } finally {
    if (showLoading) loading.value = false;
  }
};

const loadDepartments = async () => {
  try {
    const data = await apiGet('/api/departments');
    departmentsList.value = Array.isArray(data.departments) ? data.departments : [];
  } catch (error) {
    console.error('Failed to load departments', error);
    loadError.value = error.message;
    departmentsList.value = [];
  }
};

const resetNewAttendance = () => {
  newAttendance.value = { employee: '', department: '', check_in: '', check_out: '', status: 'Present' };
};

const cancelAttendanceForm = () => {
  resetNewAttendance();
  editingAttendanceId.value = null;
  showCreateAttendanceForm.value = false;
};

const startNewAttendance = () => {
  cancelAttendanceForm();
  showCreateAttendanceForm.value = true;
};

const editAttendance = (record) => {
  editingAttendanceId.value = record.id;
  newAttendance.value = {
    employee: record.employee || '',
    department: record.department || '',
    check_in: record.check_in || '',
    check_out: record.check_out || '',
    status: record.status || 'Present',
  };
  showCreateAttendanceForm.value = true;
};

const submitNewAttendance = async () => {
  try {
    if (editingAttendanceId.value !== null) {
      await apiPut(`/api/attendance/${editingAttendanceId.value}`, newAttendance.value);
    } else {
      await apiPost('/api/attendance', newAttendance.value);
    }
    await loadAttendance();
    resetNewAttendance();
    editingAttendanceId.value = null;
    showCreateAttendanceForm.value = false;
  } catch (error) {
    console.error('Failed to save attendance', error);
    alert(error.message);
  }
};

const deleteAttendance = async (record) => {
  if (!window.confirm(`Delete attendance record for ${record.employee}?`)) return;
  try {
    await apiDelete(`/api/attendance/${record.id}`);
    await loadAttendance();
  } catch (error) {
    console.error('Failed to delete attendance', error);
    alert(error.message);
  }
};

onMounted(async () => {
  loading.value = true;
  await Promise.all([loadDepartments(), loadAttendance(false)]);
  loading.value = false;
});

onBeforeUnmount(() => chartInstance?.destroy());
</script>

<style scoped>
.page-container { background: transparent; animation: fadeInView 0.4s ease-out; }
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; margin: 0 0 16px; }

/* HEADER */
.page-title-section { margin-bottom: 20px; }
.page-title-section h2 { margin: 0 0 4px; font-size: 24px; color: #0f172a; }
.breadcrumb { font-size: 13px; color: #64748b; margin: 0; }
.breadcrumb span { color: #047857; font-weight: 500; }

/* SUMMARY CARDS */
.stats-top-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 25px; }
.metric-card-box { background: white; padding: 18px; border-radius: 12px; display: flex; align-items: center; gap: 15px; border: 1px solid #e2e8f0; }
.icon-circle { width: 44px; height: 44px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 18px; }
.green-bg { background-color: #ecfdf5; }
.red-bg { background-color: #fef2f2; }
.orange-bg { background-color: #fff7ed; }
.blue-bg { background-color: #eff6ff; }
.metric-info-data h4 { margin: 0 0 4px; font-size: 12px; color: #64748b; text-transform: uppercase; letter-spacing: 0.5px; }
.metric-num-val { font-size: 24px; font-weight: 700; margin: 0 0 2px; color: #0f172a; }
.trend-indicator { font-size: 12px; font-weight: 500; }
.text-green { color: #10b981; }
.text-red, .text-orange, .text-blue { color: #047857; }

/* SPLIT LAYOUT */
.attendance-split-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 20px; }
@media (max-width: 1000px) { .attendance-split-grid { grid-template-columns: 1fr; } }

/* LEFT TABLE HOUSING */
.table-main-wrapper { min-width: 0; max-width: 100%; background: white; border-radius: 12px; border: 1px solid #e2e8f0; padding: 20px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.02); }
.controls-action-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; flex-wrap: wrap; gap: 12px; }
.filters-left-cluster { display: flex; flex: 1 1 360px; min-width: 0; gap: 10px; align-items: center; flex-wrap: wrap; }
.date-picker-mock { background: #f1f5f9; padding: 8px 14px; border-radius: 8px; font-size: 13px; font-weight: 600; color: #334155; }
.dropdown-filter { min-width: 0; max-width: 100%; min-height: 38px; padding: 8px 12px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: 13px; color: #475569; background: white; cursor: pointer; }
.action-buttons-right { display: flex; min-width: 0; max-width: 100%; align-items: center; justify-content: flex-end; gap: 8px; flex-wrap: wrap; }
.add-record-btn { background-color: #047857; color: white; border: none; padding: 8px 14px; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.filter-btn { min-height: 38px; background-color: #047857; color: white; border: none; padding: 8px 14px; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.filter-btn:hover { background-color: #059669; }
.clear-filter-btn { min-height: 38px; background: white; color: #475569; border: 1px solid #cbd5e1; padding: 8px 12px; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.clear-filter-btn:hover { background: #f8fafc; }

/* CORPORATE TABLE */
.table-scroll-area { overflow-x: auto; max-width: 100%; min-width: 0; }
.corporate-table { width: 100%; border-collapse: collapse; font-size: 13px; text-align: left; }
.corporate-table th { background-color: #f8fafc; padding: 12px 14px; font-weight: 600; color: #475569; border-bottom: 2px solid #e2e8f0; }
.corporate-table td { padding: 12px 14px; border-bottom: 1px solid #f1f5f9; color: #334155; vertical-align: middle; }
.empty-state { padding: 28px 14px !important; color: #64748b !important; text-align: center; }
.id-code { font-family: monospace; font-weight: 600; color: #64748b; }
.emp-name-td { font-weight: 600; color: #0f172a; }

/* DEPT TAGS */
.dept-label-p { padding: 2px 6px; border-radius: 4px; font-size: 11px; font-weight: 600; }
.blue-tag, .purple-tag, .orange-tag, .pink-tag { background: #d1fae5; color: #065f46; }

/* STATUS BADGES */
.status-badge-pill { padding: 3px 8px; border-radius: 6px; font-size: 11px; font-weight: 600; display: inline-block; }
.status-badge-pill.present-bg { background-color: #dcfce7; color: #16a34a; }
.status-badge-pill.late-bg, .status-badge-pill.absent-bg { background-color: #d1fae5; color: #065f46; }

/* RIGHT PANEL ANALYTICS */
.analytics-side-panel { display: flex; flex-direction: column; gap: 20px; }
.chart-card { background: white; padding: 20px; border-radius: 12px; border: 1px solid #e2e8f0; }
.chart-card h3 { margin: 0 0 15px; font-size: 14px; color: #0f172a; font-weight: 600; }
.pie-container { height: 160px; position: relative; }

/* PROGRESS BARS */
.progress-list { display: flex; flex-direction: column; gap: 12px; }
.progress-item { display: flex; flex-direction: column; gap: 4px; }
.progress-info-text { display: flex; justify-content: space-between; font-size: 12px; font-weight: 600; color: #475569; }
.bar-bg { background: #e2e8f0; height: 6px; border-radius: 4px; overflow: hidden; }
.bar-fill { height: 100%; border-radius: 4px; }

/* PAGINATION */
.table-pagination-bar { display: flex; justify-content: space-between; align-items: center; margin-top: 15px; border-top: 1px solid #f1f5f9; padding-top: 12px; }
.pagination-text-info { font-size: 12px; color: #64748b; }
.pagination-nav-cluster { display: flex; gap: 5px; }
.nav-number-btn { background: white; border: 1px solid #e2e8f0; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; border-radius: 6px; cursor: pointer; font-size: 12px; }
.nav-number-btn.active-page { background-color: #047857; color: white; border-color: #047857; }

@media (max-width: 767px) {
  .controls-action-bar { align-items: stretch; }
  .filters-left-cluster, .action-buttons-right { width: 100%; }
  .filters-left-cluster > *, .action-buttons-right > * { min-width: 0; max-width: 100%; flex: 1 1 140px; }
}

@keyframes fadeInView {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>