<template>
  <div class="dashboard-layout">
    <aside class="sidebar">
      <div class="sidebar-brand">
        <img src="/system_logo.png" alt="Logo" class="custom-sidebar-logo" />
        
        <h3>HR Portal</h3>
      </div>
      <nav class="sidebar-menu">
        <router-link to="/dashboard" class="menu-item" exact-active-class="router-link-active">📊 Dashboard</router-link>
        <template v-if="isEmployee">
          <router-link to="/dashboard/profile" class="menu-item">👤 My Profile</router-link>
          <router-link to="/dashboard/my-attendance" class="menu-item">📅 My Attendance</router-link>
          <router-link to="/dashboard/my-leaves" class="menu-item">✉️ My Leaves</router-link>
          <router-link to="/dashboard/my-payslips" class="menu-item">💵 My Payslips</router-link>
          <router-link to="/dashboard/change-password" class="menu-item">🔒 Change Password</router-link>
        </template>
        <template v-else>
          <router-link to="/dashboard/employees" class="menu-item">👥 Employees</router-link>
          <router-link to="/dashboard/departments" class="menu-item">🏢 Departments</router-link>
          <router-link to="/dashboard/attendance" class="menu-item">📅 Attendance</router-link>
          <router-link to="/dashboard/leaves" class="menu-item">✉️ Leaves</router-link>
          <router-link to="/dashboard/payroll" class="menu-item">💵 Payroll</router-link>
          <router-link to="/dashboard/reports" class="menu-item">📊 Reports</router-link>
          <router-link to="/dashboard/settings" class="menu-item">⚙️ Settings</router-link>
        </template>
      </nav>
      <div class="sidebar-footer">
        <button type="button" class="logout-btn" @click="logout">🚪 Logout</button>
      </div>
    </aside>

    <main class="main-content">
      <header class="top-header">
        <div class="header-title">
          <h2>Welcome back, {{ displayName }}! 👋</h2>
          <p class="subtitle">{{ isEmployee ? 'Your work, leave, and payslip overview.' : "Here's what's happening in your organization today." }}</p>
        </div>
        <div class="user-profile-section">
          <router-link v-if="!isEmployee" to="/dashboard/leaves" class="notification-icon" aria-label="View leave requests">🔔<span class="noti-badge">{{ pendingLeaves }}</span></router-link>
          <div class="user-avatar">{{ userInitials }}</div>
          <div class="user-info">
            <span class="user-name">{{ displayName }}</span>
            <span v-if="userPosition" class="user-role">{{ userPosition }}</span>
          </div>
        </div>
      </header>

      <div class="content-body">
        <p v-if="$route.path === '/dashboard' && dashboardError" class="api-error" role="alert">{{ dashboardError }}</p>
        <div v-if="$route.path !== '/dashboard'">
         <router-view :key="$route.fullPath" />
        </div>

        <EmployeeDashboard v-else-if="isEmployee && $route.path === '/dashboard'" />
        <div v-else>
          <EmptyState v-if="dashboardLoading" message="Loading dashboard..." />
          <div class="stats-grid">
            <div class="stat-card">
              <div class="card-icon blue-bg">👥</div>
              <div class="card-data">
                <h4>Total Employees</h4>
                <p class="stat-number">{{ totalEmployees }}</p>
              </div>
            </div>
            <div class="stat-card">
              <div class="card-icon green-bg">🏢</div>
              <div class="card-data">
                <h4>Departments</h4>
                <p class="stat-number">{{ totalDepartments }}</p>
              </div>
            </div>
            <div class="stat-card">
              <div class="card-icon orange-bg">👤</div>
              <div class="card-data">
                <h4>Present Today</h4>
                <p class="stat-number">{{ presentToday }}</p>
              </div>
            </div>
            <div class="stat-card">
              <div class="card-icon pink-bg">📅</div>
              <div class="card-data">
                <h4>Pending Leaves</h4>
                <p class="stat-number">{{ pendingLeaves }}</p>
                <router-link to="/dashboard/leaves" class="trend font-pink">View all</router-link>
              </div>
            </div>
          </div>

          <div class="dashboard-grid combo-grid">
            <div class="grid-card">
              <h3>Attendance Overview</h3>
              <div class="chart-container">
                <canvas v-if="attendanceRecords.length" id="attendanceChart" ref="attendanceChartRef"></canvas>
                <EmptyState v-else :message="dashboardLoading ? 'Loading attendance...' : 'No attendance data yet'" height="100%" />
              </div>
            </div>

            <div class="grid-card">
              <h3>Recent Leave Requests</h3>
              <div class="leave-list">
                <div v-for="leave in recentLeaveRequests" :key="leave.id" class="leave-item">
                  <div class="leave-user">{{ leave.employee }}</div>
                  <span class="leave-badge" :class="(leave.status || '').toLowerCase()">{{ leave.status || 'Unknown' }}</span>
                </div>
                <EmptyState v-if="!dashboardLoading && !recentLeaveRequests.length" message="No leave requests yet" />
              </div>
            </div>
          </div>

          <div class="dashboard-grid combo-grid bottom-grid">
            <div class="grid-card table-card">
              <h3>Recent Employees</h3>
              <table class="minimal-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Department</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="employee in recentEmployees" :key="employee.ID || employee.employee_id">
                    <td>{{ employee.employee_id }}</td>
                    <td>{{ employee.full_name }}</td>
                    <td>{{ employee.department }}</td>
                    <td><span class="status-t" :class="(employee.status || '').toLowerCase() === 'active' ? 'active' : ''">{{ employee.status || 'Unknown' }}</span></td>
                  </tr>
                  <tr v-if="!dashboardLoading && !recentEmployees.length">
                    <td colspan="4"><EmptyState message="No employees yet" /></td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="grid-card">
              <h3>Department Distribution</h3>
              <div class="chart-container pie-size">
                <canvas v-if="departmentDistribution.length" id="deptPieChart" ref="deptPieChartRef"></canvas>
                <EmptyState v-else :message="dashboardLoading ? 'Loading department data...' : 'No department data yet'" height="100%" />
              </div>
            </div>
          </div>

        </div> </div>
    </main>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, nextTick, ref, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Chart from 'chart.js/auto';
import { apiGet, clearAuthSession } from '@/api';
import EmptyState from '@/components/EmptyState.vue';
import EmployeeDashboard from './EmployeeDashboard.vue';

const route = useRoute();
const router = useRouter();
const attendanceChartRef = ref(null);
const deptPieChartRef = ref(null);
const totalEmployees = ref(0);
const totalDepartments = ref(0);
const presentToday = ref(0);
const pendingLeaves = ref(0);
const dashboardLoading = ref(false);
const dashboardError = ref('');
const authUser = ref(readStoredUser());
const isEmployee = computed(() => authUser.value.role === 'Employee');
const employeeRecords = ref([]);
const leaveRequests = ref([]);
const attendanceRecords = ref([]);
const displayName = computed(() => authUser.value.full_name || authUser.value.email || 'User');
const userPosition = computed(() => authUser.value.position || '');
const userInitials = computed(() => displayName.value.split(/\s+/).map(part => part[0]).join('').slice(0, 2).toUpperCase());
const recentEmployees = computed(() => [...employeeRecords.value]
  .sort((a, b) => {
    const dateA = new Date(a.CreatedAt).getTime();
    const dateB = new Date(b.CreatedAt).getTime();
    if (Number.isNaN(dateA)) return Number.isNaN(dateB) ? 0 : 1;
    if (Number.isNaN(dateB)) return -1;
    return dateB - dateA;
  })
  .slice(0, 5));
const recentLeaveRequests = computed(() => leaveRequests.value.slice(0, 5));
const departmentDistribution = computed(() => {
  const counts = new Map();
  for (const employee of employeeRecords.value) {
    const department = employee.department?.trim();
    if (department) counts.set(department, (counts.get(department) || 0) + 1);
  }
  return [...counts].map(([label, count]) => ({ label, count }));
});
const chartInstances = {
  attendance: null,
  deptPie: null
};
function readStoredUser() {
  try {
    return JSON.parse(localStorage.getItem('authUser') || '{}');
  } catch {
    return {};
  }
}
const updateAuthUser = (event) => {
  authUser.value = event.detail;
};

const loadDashboard = async () => {
  if (isEmployee.value) return;
  dashboardLoading.value = true;
  dashboardError.value = '';
  try {
    const [summary, employees, leaves, attendance] = await Promise.all([
      apiGet('/api/dashboard-summary'),
      apiGet('/api/employees'),
      apiGet('/api/leaves'),
      apiGet('/api/attendance'),
    ]);
    totalEmployees.value = Number(summary.total_employees) || 0;
    totalDepartments.value = Number(summary.total_departments) || 0;
    presentToday.value = Number(summary.present_today) || 0;
    pendingLeaves.value = Number(summary.pending_leaves) || 0;
    employeeRecords.value = Array.isArray(employees) ? employees : [];
    leaveRequests.value = Array.isArray(leaves.requests) ? leaves.requests : [];
    attendanceRecords.value = Array.isArray(attendance.attendance) ? attendance.attendance : [];
    await createCharts();
  } catch (error) {
    console.error('Unable to load dashboard summary', error);
    dashboardError.value = `Unable to load dashboard data: ${error.message}`;
    totalEmployees.value = 0;
    totalDepartments.value = 0;
    presentToday.value = 0;
    pendingLeaves.value = 0;
    employeeRecords.value = [];
    leaveRequests.value = [];
    attendanceRecords.value = [];
    await createCharts();
  } finally {
    dashboardLoading.value = false;
  }
};

const logout = () => {
  clearAuthSession();
  router.push('/login');
};

const destroyCharts = () => {
  if (chartInstances.attendance) {
    chartInstances.attendance.destroy();
    chartInstances.attendance = null;
  }
  if (chartInstances.deptPie) {
    chartInstances.deptPie.destroy();
    chartInstances.deptPie = null;
  }
};

const createCharts = async () => {
  await nextTick();
  await nextTick();
  destroyCharts();

  const attendanceLabels = [...new Set(attendanceRecords.value.map(record => record.status || 'Unknown'))];
  const ctxLine = attendanceChartRef.value;
  if (ctxLine && attendanceRecords.value.length) {
    chartInstances.attendance = new Chart(ctxLine, {
      type: 'bar',
      data: {
        labels: attendanceLabels,
        datasets: [{
          label: 'Attendance records',
          data: attendanceLabels.map(status => attendanceRecords.value.filter(record => (record.status || 'Unknown') === status).length),
          borderColor: '#10b981',
          backgroundColor: 'rgba(16, 185, 129, 0.35)',
          borderWidth: 2,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } }
      }
    });
  }

  const ctxPie = deptPieChartRef.value;
  if (ctxPie && departmentDistribution.value.length) {
    chartInstances.deptPie = new Chart(ctxPie, {
      type: 'doughnut',
      data: {
        labels: departmentDistribution.value.map(department => department.label),
        datasets: [{
          data: departmentDistribution.value.map(department => department.count),
          backgroundColor: ['#244b3e', '#047857', '#10b981', '#34d399', '#a7f3d0']
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { position: 'right' } }
      }
    });
  }
};

onMounted(async () => {
  window.addEventListener('auth-user-updated', updateAuthUser);
  if (route.path === '/dashboard') {
    await loadDashboard();
  }
});

watch(
  () => route.fullPath,
  (newPath) => {
    if (newPath === '/dashboard') {
      loadDashboard();
    } else {
      destroyCharts();
    }
  }
);

onBeforeUnmount(() => {
  window.removeEventListener('auth-user-updated', updateAuthUser);
  destroyCharts();
});
</script>

<style scoped>


/* GENERAL SKELETON */
.dashboard-layout { 
  min-height: 100vh;
  height: 100vh;
  overflow: hidden;
  background-color: #f4f7f6; /* පිරිසිදු Off-white පසුබිමක් */
  font-family: 'Inter', 'Segoe UI', sans-serif; 
  color: #0f2922; /* තද කොළ පැහැයට හුරු Text කලර් එකක් */
}
.main-content { 
  height: 100vh;
  margin-left: 260px;
  width: calc(100% - 260px);
  box-sizing: border-box;
  display: flex; 
  flex-direction: column; 
  overflow-x: hidden;
  overflow-y: auto;
}
.content-body { 
  min-width: 0;
  padding: 25px; 
}
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; margin: 0 0 16px; }

/* 1. DEEP EMERALD SIDEBAR */
.sidebar { 
  position: fixed;
  top: 0;
  left: 0;
  z-index: 20;
  width: 260px;
  height: 100vh;
  box-sizing: border-box;
  background-color: #386657; 
  color: #e6f4f0; 
  display: flex; 
  flex-direction: column; 
  box-shadow: 4px 0 15px rgba(6, 44, 34, 0.05);
 
  border-radius: 20px  !important; 
  
 
  overflow-x: hidden !important;
  overflow-y: auto !important;
}
.sidebar-brand {
  display: flex;
  align-items: center;
  justify-content: flex-start; 
  gap: 10px;
  flex: 0 0 64px;
  height: 64px;
  padding: 8px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.custom-sidebar-logo {
  flex: 0 0 36px;
  width: 36px;
  height: 36px;
  object-fit: contain;
}
.logo-icon { 
  background: #10b981; /* Mint / Vibrant Green */
  width: 30px; 
  height: 30px; 
  border-radius: 8px; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  font-weight: bold; 
  font-size: 14px;
  color: #062c22;
}
.sidebar-brand h3 { 
  font-size: 16px; 
  margin: 0; 
  white-space: nowrap;
  color: #ffffff; 
  font-weight: 600;
}
.sidebar-menu { 
  flex: 1; 
  padding: 25px 0; 
}
.menu-item {
  display: flex;
  align-items: center;
  padding: 14px 24px;
  
  
  color: #ffffff !important; 
  opacity: 0.85 !important; 
  
  text-decoration: none;
  font-weight: 500;
  transition: 0.2s ease;
  border-left: 4px solid transparent;
}

.menu-item:hover {
  background-color: rgba(230, 244, 240, 0.05);
  

  color: #ffffff !important;
  opacity: 1 !important;
}


.menu-item.router-link-active {

  background: rgba(16, 185, 129, 0.1) !important; 
  
 
  color: #10b981 !important;
  border-left: 4px solid #10b981 !important;
  opacity: 1 !important;
  font-weight: 600;
}
.sidebar-footer { 
  padding: 20px; 
  border-top: 1px solid rgba(230, 244, 240, 0.1); 
}
.logout-btn { 
  color: #d1fae5; 
  text-decoration: none; 
  font-weight: 500; 
  background: transparent;
  border: 0;
  padding: 0;
  font: inherit;
  cursor: pointer;
}

/* 2. TOP HEADER STYLE */
.top-header { 
  background: white; 
  height: 64px;
  min-height: 64px;
  padding: 0 30px;
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  border-bottom: 1px solid #d7e5df;
  box-shadow: 0 1px 3px rgba(24, 53, 43, 0.08);
}
.top-header h2 { 
  margin: 0; 
  font-size: 22px; 
  color: #062c22; 
  font-weight: 700;
}
.subtitle, .breadcrumb { 
  margin: 4px 0 0; 
  font-size: 13px; 
  color: #52796f; 
}
.breadcrumb .active-crumb {
  color: #10b981;
  font-weight: 600;
}
.user-profile-section { 
  display: flex; 
  align-items: center; 
  gap: 15px; 
}
.notification-icon { 
  font-size: 20px; 
  position: relative; 
  cursor: pointer; 
  color: #52796f; 
  text-decoration: none;
}
.noti-badge { 
  position: absolute; 
  top: -5px; 
  right: -5px; 
  background: #047857; 
  color: white; 
  font-size: 10px; 
  padding: 2px 5px; 
  border-radius: 50%; 
}
.user-avatar { 
  background: #244b3e; 
  color: #6ee7b7; 
  width: 38px; 
  height: 38px; 
  border-radius: 50%; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  font-weight: bold; 
}
.user-info { 
  display: flex; 
  flex-direction: column; 
}
.user-name { 
  font-weight: 600; 
  font-size: 14px; 
  color: #062c22;
}
.user-role { 
  font-size: 11px; 
  color: #52796f; 
}

/* 3. PREMIUM SUMMARY CARDS */
.stats-grid { 
  display: grid; 
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); 
  gap: 20px; 
  margin-bottom: 25px; 
}
.stat-card { 
  background: white; 
  padding: 20px; 
  border-radius: 16px; 
  display: flex; 
  align-items: center; 
  gap: 15px; 
  box-shadow: 0 4px 10px rgba(6, 44, 34, 0.02); 
  border: 1px solid #e5eef0;
  transition: transform 0.2s; 
}
.stat-card:hover { 
  transform: translateY(-2px); 
  box-shadow: 0 6px 15px rgba(6, 44, 34, 0.04);
}
.card-icon { 
  font-size: 22px; 
  width: 45px; 
  height: 45px; 
  border-radius: 50%; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
}
.blue-bg { background: #e6f4f0; color: #062c22; }
.green-bg { background: #ecfdf5; color: #10b981; }
.orange-bg { background: #d1fae5; color: #065f46; }
.pink-bg { background: #fff1f2; color: #f43f5e; }

.card-data h4 { 
  margin: 0 0 4px; 
  font-size: 12px; 
  color: #52796f; 
  font-weight: 600; 
  text-transform: uppercase; 
  letter-spacing: 0.5px; 
}
.stat-number { 
  font-size: 24px; 
  font-weight: 700; 
  margin: 0; 
  color: #062c22;
}
.trend { 
  font-size: 12px; 
}
.font-green { color: #10b981; font-weight: 500; }
.font-orange { color: #047857; font-weight: 500; }
.font-pink { color: #f43f5e; font-weight: 500; cursor: pointer; }
.font-gray { color: #94a3b8; }

/* 4. CHARTS SECTION LAYOUT */
.combo-grid { 
  display: grid; 
  grid-template-columns: 2fr 1fr; 
  gap: 20px; 
  margin-bottom: 25px; 
}
@media (max-width: 850px) { 
  .combo-grid { grid-template-columns: 1fr; } 
}
@media (max-width: 700px) {
  .sidebar { width: 210px; }
  .main-content { margin-left: 210px; width: calc(100% - 210px); }
  .content-body { padding: 16px; }
  .top-header { padding: 0 16px; }
}
.grid-card { 
  background: white; 
  padding: 20px; 
  border-radius: 16px; 
  border: 1px solid #e5eef0; 
  box-shadow: 0 4px 10px rgba(6, 44, 34, 0.02);
}
.grid-card h3 { 
  margin: 0 0 15px; 
  font-size: 14px; 
  color: #062c22; 
  font-weight: 600;
}
.chart-container { 
  height: 220px; 
  position: relative; 
}
.pie-size { 
  height: 180px; 
}

/* 5. LEAVE REQUESTS & MINIMAL TABLES */
.leave-list { 
  display: flex; 
  flex-direction: column; 
  gap: 12px; 
}
.leave-item { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  padding: 10px 14px; 
  background: #f8fafc; 
  border-radius: 10px; 
  border: 1px solid #e5eef0;
}
.leave-user { 
  font-size: 13px; 
  font-weight: 600; 
  color: #062c22;
}
.leave-badge { 
  font-size: 11px; 
  padding: 3px 8px; 
  border-radius: 12px; 
  font-weight: 600; 
}
.leave-badge.pending { background: #d1fae5; color: #065f46; }
.leave-badge.approved { background: #e6f4f0; color: #062c22; }

/* TABLES STYLING */
.minimal-table { 
  width: 100%; 
  border-collapse: collapse; 
  font-size: 13px; 
}
.minimal-table th, .minimal-table td { 
  padding: 12px 14px; 
  text-align: left; 
  border-bottom: 1px solid #e5eef0; 
}
.minimal-table th { 
  background: #f8fafc; 
  color: #52796f; 
  font-weight: 600; 
}
.status-t { 
  padding: 2px 6px; 
  border-radius: 4px; 
  font-size: 11px; 
  font-weight: bold; 
}
.status-t.active { 
  background: #ecfdf5; 
  color: #10b981; 
}
</style>