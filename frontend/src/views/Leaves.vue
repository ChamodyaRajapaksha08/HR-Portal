<template>
  <div class="page-container">
    <p v-if="loadError" class="api-error" role="alert">{{ loadError }}</p>
    <div class="page-title-section">
      <h2>Leaves</h2>
      <p class="breadcrumb">Dashboard > <span>Leaves</span></p>
    </div>

    <div class="stats-top-grid">
      <div class="metric-card-box">
        <div class="icon-circle purple-bg">✉️</div>
        <div class="metric-info-data">
          <h4>Total Leave Requests</h4>
          <p class="metric-num-val">{{ summary.total_leave_requests }}</p>
          <span class="trend-indicator text-green">▲ 3 this month</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle orange-bg">⏳</div>
        <div class="metric-info-data">
          <h4>Pending Requests</h4>
          <p class="metric-num-val">{{ summary.pending_requests }}</p>
          <span class="trend-indicator text-orange">● 40.00%</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle green-bg">✅</div>
        <div class="metric-info-data">
          <h4>Approved</h4>
          <p class="metric-num-val">{{ summary.approved }}</p>
          <span class="trend-indicator text-green">● 46.67%</span>
        </div>
      </div>
      <div class="metric-card-box">
        <div class="icon-circle red-bg">❌</div>
        <div class="metric-info-data">
          <h4>Rejected</h4>
          <p class="metric-num-val">{{ summary.rejected }}</p>
          <span class="trend-indicator text-red">● 13.33%</span>
        </div>
      </div>
    </div>

    <div class="leaves-split-grid">
      
      <div ref="tableSection" class="table-main-wrapper">
        <div class="table-tab-filters">
          <div class="filters-right-cluster">
            <div class="search-field-container">
              <span class="search-lens">🔍</span>
              <input v-model="searchQuery" type="search" placeholder="Search leave requests..." class="search-input-element" />
            </div>
            <select v-model="selectedStatus" class="dropdown-filter">
              <option v-for="status in leaveStatuses" :key="status" :value="status">{{ status }}</option>
            </select>
            <select v-model="selectedLeaveType" class="dropdown-filter">
              <option value="">All Leave Types</option>
              <option v-for="leaveType in leaveTypes" :key="leaveType" :value="leaveType">{{ leaveType }}</option>
            </select>
            <button class="add-record-btn" @click="startNewLeave">+ Add Request</button>
          </div>
        </div>

        <form v-if="showCreateLeaveForm" class="create-record-form" @submit.prevent="submitNewLeaveRequest">
          <h3 class="leave-form-title">{{ editingLeaveId ? 'Edit Leave Request' : 'New Leave Request' }}</h3>
          <div class="form-grid">
            <label>
              Employee
              <input v-model="newLeave.employee" type="text" placeholder="Kasun Perera" required />
            </label>
            <label>
              Leave Type
              <input v-model="newLeave.leave_type" type="text" placeholder="Annual Leave" required />
            </label>
            <label>
              From Date
              <input v-model="newLeave.from_date" type="date" required />
            </label>
            <label>
              To Date
              <input v-model="newLeave.to_date" type="date" required />
            </label>
            <label>
              Days
              <input v-model.number="newLeave.days" type="number" min="1" placeholder="4" required />
            </label>
            <label>
              Reason
              <input v-model="newLeave.reason" type="text" placeholder="Reason for leave" />
            </label>
            <label>
              Status
              <select v-model="newLeave.status">
                <option>Pending</option>
                <option>Approved</option>
                <option>Rejected</option>
              </select>
            </label>
          </div>
          <div class="form-actions">
            <button type="submit" class="submit-btn">{{ editingLeaveId ? 'Update Request' : 'Save Request' }}</button>
            <button type="button" class="cancel-btn" @click="cancelLeaveForm">Cancel</button>
          </div>
        </form>

        <div class="table-scroll-area">
          <table class="corporate-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Employee</th>
                <th>Leave Type</th>
                <th>Date Range</th>
                <th>Days</th>
                <th>Reason</th>
                <th>Status</th>
                <th>Applied On</th>
                <th style="text-align: center;">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="request in pagedLeaveRequests" :key="request.id">
                <td class="id-code">{{ request.id }}</td>
                <td class="emp-name-td">{{ request.employee }}</td>
                <td>{{ request.leave_type }}</td>
                <td>{{ request.from_date }} to {{ request.to_date }}</td>
                <td>{{ request.days }}</td>
                <td class="reason-cell">{{ request.reason }}</td>
                <td><span class="status-badge" :class="statusClass(request.status)">{{ request.status }}</span></td>
                <td class="date-applied">{{ request.applied_on }}</td>
                <td>
                  <div class="action-btn-group">
                    <button class="table-icon-btn edit-color" title="Edit leave request" @click="editLeave(request)">📝</button>
                    <button class="table-icon-btn delete-color" title="Delete leave request" @click="deleteLeave(request)">🗑️</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="table-pagination-bar">
          <span class="pagination-text-info">Showing {{ paginationStart }}–{{ paginationEnd }} of {{ filteredLeaveRequests.length }} requests</span>
          <div class="pagination-nav-cluster">
            <button class="nav-arrow-btn" :disabled="currentPage === 1" @click="currentPage--">‹</button>
            <button v-for="page in pageCount" :key="page" class="nav-number-btn" :class="{ 'active-page': page === currentPage }" @click="currentPage = page">{{ page }}</button>
            <button class="nav-arrow-btn" :disabled="currentPage >= pageCount" @click="currentPage++">›</button>
          </div>
        </div>
      </div>

      <div class="analytics-side-panel">
        <div class="chart-card">
          <h3>Leave Type Distribution</h3>
          <div class="pie-container">
            <canvas id="leaveTypeDoughnutChart"></canvas>
          </div>
        </div>

        <div class="chart-card upcoming-card">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <h3>Upcoming Leaves</h3>
            <button type="button" class="view-all-link" @click="showAllLeaves">View all</button>
          </div>
          <div class="upcoming-list">
            <div v-for="(leave, index) in upcomingLeaves" :key="leave.id" class="upcoming-item">
              <div class="upcoming-user-info">
                <span class="u-name">{{ leave.employee }}</span>
                <span class="u-type">{{ leave.leave_type }}</span>
              </div>
              <div class="upcoming-date-info">
                <span class="u-date">{{ formatLeaveDate(leave.from_date) }} to {{ formatLeaveDate(leave.to_date) }}</span>
                <span class="days-badge" :class="index % 2 === 0 ? 'b-blue' : 'b-purple'">{{ leave.days }} {{ leave.days === 1 ? 'day' : 'days' }}</span>
              </div>
            </div>
            <p v-if="!upcomingLeaves.length" class="no-upcoming-leaves">No upcoming approved or pending leaves.</p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onBeforeUnmount, ref, watch } from 'vue';
import Chart from 'chart.js/auto';
import { apiDelete, apiGet, apiPost, apiPut } from '@/api';

const summary = ref({ total_leave_requests: 0, pending_requests: 0, approved: 0, rejected: 0 });
const leaveRequests = ref([]);
const tableSection = ref(null);
const showCreateLeaveForm = ref(false);
const newLeave = ref({ employee: '', leave_type: '', from_date: '', to_date: '', days: 1, reason: '', status: 'Pending', applied_on: '' });
const searchQuery = ref('');
const selectedStatus = ref('All Requests');
const selectedLeaveType = ref('');
const editingLeaveId = ref(null);
const loadError = ref('');
const currentPage = ref(1);
const pageSize = 8;
const leaveStatuses = ['All Requests', 'Pending', 'Approved', 'Rejected'];
const leaveTypes = ['Annual', 'Medical', 'Personal', 'Sick'];
const upcomingLeaves = computed(() => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return leaveRequests.value
    .filter(request => ['Approved', 'Pending'].includes(request.status) && parseLeaveDate(request.from_date) >= today)
    .sort((a, b) => parseLeaveDate(a.from_date) - parseLeaveDate(b.from_date))
    .slice(0, 4);
});
const filteredLeaveRequests = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  return leaveRequests.value.filter(request => {
    const matchesQuery = !query || [request.employee, request.leave_type, request.reason, request.status]
      .some(value => (value || '').toLowerCase().includes(query));
    const matchesStatus = !selectedStatus.value || selectedStatus.value === 'All Requests' || request.status === selectedStatus.value;
    const matchesType = !selectedLeaveType.value || request.leave_type === selectedLeaveType.value;
    return matchesQuery && matchesStatus && matchesType;
  });
});
const pageCount = computed(() => Math.max(1, Math.ceil(filteredLeaveRequests.value.length / pageSize)));
const pagedLeaveRequests = computed(() => filteredLeaveRequests.value.slice((currentPage.value - 1) * pageSize, currentPage.value * pageSize));
const paginationStart = computed(() => filteredLeaveRequests.value.length ? (currentPage.value - 1) * pageSize + 1 : 0);
const paginationEnd = computed(() => Math.min(currentPage.value * pageSize, filteredLeaveRequests.value.length));
watch([searchQuery, selectedStatus, selectedLeaveType], () => { currentPage.value = 1; });
let leaveChart = null;

const parseLeaveDate = (date) => {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(date || '')) return new Date(NaN);
  const [year, month, day] = date.split('-').map(Number);
  const parsedDate = new Date(year, month - 1, day);
  return parsedDate.getFullYear() === year && parsedDate.getMonth() === month - 1 && parsedDate.getDate() === day
    ? parsedDate
    : new Date(NaN);
};

const formatLeaveDate = (date) => {
  const parsedDate = parseLeaveDate(date);
  return Number.isNaN(parsedDate.getTime())
    ? date
    : parsedDate.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' });
};

const showAllLeaves = async () => {
  searchQuery.value = '';
  selectedStatus.value = 'All Requests';
  selectedLeaveType.value = '';
  currentPage.value = 1;
  await nextTick();
  tableSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' });
};

const statusClass = (status) => {
  if (status === 'Approved') return 'approved';
  if (status === 'Pending') return 'pending';
  if (status === 'Rejected') return 'rejected';
  return '';
};

const loadLeaves = async () => {
  try {
    const data = await apiGet('/api/leaves');
    summary.value = data.summary || summary.value;
    leaveRequests.value = data.requests || [];
  } catch (error) {
    console.error('Failed to load leaves', error);
    loadError.value = error.message;
  }
};

const resetNewLeaveRequest = () => {
  newLeave.value = { employee: '', leave_type: '', from_date: '', to_date: '', days: 1, reason: '', status: 'Pending', applied_on: '' };
};

const cancelLeaveForm = () => {
  resetNewLeaveRequest();
  editingLeaveId.value = null;
  showCreateLeaveForm.value = false;
};

const startNewLeave = () => {
  cancelLeaveForm();
  showCreateLeaveForm.value = true;
};

const editLeave = (request) => {
  editingLeaveId.value = request.id;
  newLeave.value = {
    employee: request.employee || '', leave_type: request.leave_type || '',
    from_date: request.from_date || '', to_date: request.to_date || '',
    days: Number(request.days) || 1, reason: request.reason || '',
    status: request.status || 'Pending', applied_on: request.applied_on || '',
  };
  showCreateLeaveForm.value = true;
};

const submitNewLeaveRequest = async () => {
  try {
    if (editingLeaveId.value !== null) {
      await apiPut(`/api/leaves/${editingLeaveId.value}`, newLeave.value);
    } else {
      await apiPost('/api/leaves', newLeave.value);
    }
    await loadLeaves();
    resetNewLeaveRequest();
    editingLeaveId.value = null;
    showCreateLeaveForm.value = false;
  } catch (error) {
    console.error('Failed to save leave request', error);
    alert(error.message);
  }
};

const deleteLeave = async (request) => {
  if (!window.confirm(`Delete leave request for ${request.employee}?`)) return;
  try {
    await apiDelete(`/api/leaves/${request.id}`);
    await loadLeaves();
  } catch (error) {
    console.error('Failed to delete leave request', error);
    alert(error.message);
  }
};

onMounted(async () => {
  await loadLeaves();
  const ctx = document.getElementById('leaveTypeDoughnutChart');
  if (ctx) {
    leaveChart = new Chart(ctx, {
      type: 'doughnut',
      data: {
        labels: leaveTypes,
        datasets: [{
          data: [46.7, 26.7, 20.0, 6.7],
          backgroundColor: ['#244b3e', '#047857', '#10b981', '#6ee7b7'],
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { position: 'right', labels: { boxWidth: 10, font: { size: 11 } } }
        }
      }
    });

    onBeforeUnmount(() => leaveChart?.destroy());
  }
});
</script>

<style scoped>
.page-container { background: transparent; animation: fadeInView 0.4s ease-out; width: 100%; box-sizing: border-box; }
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; margin: 0 0 16px; }

/* HEADER */
.page-title-section { margin-bottom: 20px; }
.page-title-section h2 { margin: 0 0 4px; font-size: 24px; color: #0f172a; }
.breadcrumb { font-size: 13px; color: #64748b; margin: 0; }
.breadcrumb span { color: #047857; font-weight: 500; }

/* TOP CARDS GRID */
.stats-top-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 20px; margin-bottom: 25px; }
.metric-card-box { background: white; padding: 18px; border-radius: 12px; display: flex; align-items: center; gap: 15px; border: 1px solid #e2e8f0; }
.icon-circle { width: 44px; height: 44px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 18px; }
.purple-bg, .orange-bg, .green-bg, .red-bg { background-color: #d1fae5; color: #065f46; }
.metric-info-data h4 { margin: 0 0 4px; font-size: 12px; color: #64748b; text-transform: uppercase; letter-spacing: 0.5px; }
.metric-num-val { font-size: 24px; font-weight: 700; margin: 0 0 2px; color: #0f172a; }
.trend-indicator { font-size: 12px; font-weight: 500; }
.text-green { color: #10b981; }
.text-orange, .text-red { color: #047857; }

/* SPLIT GRID RESOLUTION (Prevents overlapping) */
.leaves-split-grid { display: grid; grid-template-columns: minmax(0, 2.7fr) minmax(300px, 1fr); gap: 20px; align-items: start; }
@media (max-width: 1200px) { .leaves-split-grid { grid-template-columns: minmax(0, 1fr); } }

/* LEFT TABLE WRAPPER */
.table-main-wrapper { background: white; border-radius: 12px; border: 1px solid #e2e8f0; padding: 20px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.02); min-width: 0; }
.table-tab-filters { display: flex; align-items: center; margin-bottom: 20px; }

.filters-right-cluster { display: grid; grid-template-columns: minmax(180px, 1fr) repeat(2, minmax(135px, 155px)) auto; gap: 12px; align-items: center; width: 100%; }
.search-field-container { position: relative; width: 100%; min-width: 0; }
.search-lens { position: absolute; left: 10px; top: 50%; transform: translateY(-50%); color: #94a3b8; font-size: 12px; }
.search-input-element { width: 100%; box-sizing: border-box; padding: 9px 10px 9px 30px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: 13px; background-color: #f8fafc; outline: none; }
.search-input-element:focus, .dropdown-filter:focus { border-color: #10b981; outline: 2px solid rgba(16, 185, 129, 0.12); }
.dropdown-filter { width: 100%; min-width: 0; box-sizing: border-box; padding: 9px 10px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: 13px; color: #475569; background: white; cursor: pointer; }
.add-record-btn { min-height: 37px; padding: 8px 14px; border: 0; border-radius: 8px; color: white; background: #047857; font-size: 13px; font-weight: 600; white-space: nowrap; cursor: pointer; transition: background-color 0.2s; }
.add-record-btn:hover { background: #059669; }
.filter-btn { background-color: #047857; color: white; border: none; padding: 8px 14px; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }

@media (max-width: 700px) {
  .filters-right-cluster { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .search-field-container { grid-column: 1 / -1; }
  .add-record-btn { width: 100%; }
}

/* FORM CARD */
.create-record-form { display: block; position: relative; width: 100%; box-sizing: border-box; margin: 0 0 24px; padding: 20px; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; }
.leave-form-title { margin: 0 0 16px; color: #0f172a; font-size: 16px; font-weight: 600; }
.form-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 16px; margin-bottom: 18px; }
.form-grid label { display: flex; flex-direction: column; gap: 7px; color: #475569; font-size: 13px; font-weight: 600; }
.form-grid input, .form-grid select { width: 100%; box-sizing: border-box; min-height: 40px; padding: 9px 11px; color: #0f172a; background: white; border: 1px solid #cbd5e1; border-radius: 8px; font: inherit; font-weight: 400; }
.form-grid input:focus, .form-grid select:focus { outline: 2px solid rgba(16, 185, 129, 0.2); border-color: #10b981; }
.form-actions { display: flex; gap: 10px; flex-wrap: wrap; }
.submit-btn, .cancel-btn { padding: 10px 16px; border: 1px solid transparent; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.submit-btn { color: white; background: #10b981; }
.submit-btn:hover { background: #059669; }
.cancel-btn { color: #475569; background: white; border-color: #cbd5e1; }
.cancel-btn:hover { background: #f1f5f9; }

/* CORPORATE TABLE WITH EXTRA OVERFLOW CONSTRAINTS */
.table-scroll-area { overflow-x: auto; width: 100%; -webkit-overflow-scrolling: touch; }
.corporate-table { width: 100%; border-collapse: collapse; font-size: 13px; text-align: left; white-space: nowrap; }
.corporate-table th { background-color: #f8fafc; padding: 12px 14px; font-weight: 600; color: #475569; border-bottom: 2px solid #e2e8f0; }
.corporate-table td { padding: 12px 14px; border-bottom: 1px solid #f1f5f9; color: #334155; vertical-align: middle; }
.id-code { font-family: monospace; font-weight: 600; color: #64748b; }
.emp-name-td { font-weight: 600; color: #0f172a; }
.reason-cell { max-width: 150px; overflow: hidden; text-overflow: ellipsis; }

/* STATUS BADGES */
.status-badge { padding: 3px 8px; border-radius: 6px; font-size: 11px; font-weight: 600; display: inline-block; }
.status-badge.pending { background-color: #d1fae5; color: #065f46; }
.status-badge.approved { background-color: #dcfce7; color: #16a34a; }
.status-badge.rejected { background-color: #d1fae5; color: #065f46; }

/* ACTIONS */
.action-btn-group { display: flex; gap: 6px; justify-content: center; }
.table-icon-btn { background: none; border: 1px solid #e2e8f0; padding: 5px 8px; border-radius: 6px; cursor: pointer; font-size: 13px; }
.table-icon-btn:hover { background-color: #f8fafc; }
.edit-color { color: #047857; }
.delete-color { color: #047857; }

/* RIGHT PANEL ANALYTICS (Fixed structural widths) */
.analytics-side-panel { display: flex; flex-direction: column; gap: 20px; min-width: 300px; }
.chart-card { background: white; padding: 20px; border-radius: 12px; border: 1px solid #e2e8f0; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.02); }
.chart-card h3 { margin: 0; font-size: 14px; color: #0f172a; font-weight: 600; }
.pie-container { height: 160px; position: relative; margin-top: 15px; }
.view-all-link { padding: 0; border: 0; background: transparent; font: inherit; font-size: 12px; color: #047857; font-weight: 600; cursor: pointer; }
.view-all-link:hover { color: #059669; text-decoration: underline; }

/* UPCOMING LEAVES ITEMS */
.upcoming-list { display: flex; flex-direction: column; gap: 12px; margin-top: 15px; }
.upcoming-item { padding: 10px; background: #f8fafc; border-radius: 8px; border: 1px solid #f1f5f9; display: flex; flex-direction: column; gap: 4px; }
.upcoming-user-info { display: flex; justify-content: space-between; font-size: 13px; font-weight: 600; color: #18352b; }
.upcoming-user-info .u-type { font-size: 11px; color: #64748b; font-weight: 500; }
.upcoming-date-info { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: #64748b; }
.no-upcoming-leaves { margin: 0; padding: 12px; border-radius: 8px; background: #f8fafc; color: #64748b; font-size: 12px; text-align: center; }
.days-badge { font-size: 11px; padding: 2px 6px; border-radius: 4px; font-weight: 600; }
.days-badge.b-blue, .days-badge.b-purple { background: #d1fae5; color: #065f46; }

/* PAGINATION */
.table-pagination-bar { display: flex; justify-content: space-between; align-items: center; margin-top: 15px; border-top: 1px solid #f1f5f9; padding-top: 12px; }
.pagination-text-info { font-size: 12px; color: #64748b; }
.pagination-nav-cluster { display: flex; gap: 5px; }
.nav-number-btn { background: white; border: 1px solid #e2e8f0; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; border-radius: 6px; cursor: pointer; font-size: 12px; }
.nav-number-btn.active-page { background-color: #047857; color: white; border-color: #047857; }

@keyframes fadeInView {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>