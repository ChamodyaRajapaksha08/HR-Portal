<template>
  <div class="dashboard-layout">
    <aside class="sidebar">
      <div class="sidebar-brand">
        <img src="/system_logo.png" alt="Logo" class="custom-sidebar-logo" />
        
        <h3>HR Portal</h3>
      </div>
      <nav class="sidebar-menu">
        <router-link to="/dashboard" class="menu-item" exact-active-class="router-link-active">📊 Dashboard</router-link>

        <router-link to="/dashboard/employees" class="menu-item">👥 Employees</router-link>
        <router-link to="/dashboard/departments" class="menu-item">🏢 Departments</router-link>
        <router-link to="/dashboard/attendance" class="menu-item">📅 Attendance</router-link>
        <router-link to="/dashboard/leaves" class="menu-item">✉️ Leaves</router-link>

        <router-link to="/dashboard/payroll" class="menu-item">💵 Payroll</router-link>
        <router-link to="/dashboard/reports" class="menu-item">📊 Reports</router-link>
        <router-link to="/dashboard/settings" class="menu-item">⚙️ Settings</router-link>
      </nav>
      <div class="sidebar-footer">
        <button type="button" class="logout-btn" @click="logout">🚪 Logout</button>
      </div>
    </aside>

    <main class="main-content">
      <header class="top-header">
        <div class="header-title">
          <h2>Welcome back, {{ displayName }}! 👋</h2>
          <p class="subtitle">Here's what's happening in your organization today.</p>
        </div>
        <div class="user-profile-section">
          <router-link to="/dashboard/leaves" class="notification-icon" aria-label="View leave requests">🔔<span class="noti-badge">{{ pendingLeaves }}</span></router-link>
          <div class="user-avatar">AD</div>
          <div class="user-info">
            <span class="user-name">{{ displayName }}</span>
            <span class="user-role">{{ userPosition }}</span>
          </div>
        </div>
      </header>

      <div class="content-body">
        <p v-if="$route.path === '/dashboard' && dashboardError" class="api-error" role="alert">{{ dashboardError }}</p>
        <div v-if="$route.path !== '/dashboard'">
         <router-view :key="$route.fullPath" />
        </div>

        <div v-else>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="card-icon blue-bg">👥</div>
              <div class="card-data">
                <h4>Total Employees</h4>
                <p class="stat-number">{{ totalEmployees }}</p>
                <span class="trend font-green">▲ 8 this month</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="card-icon green-bg">🏢</div>
              <div class="card-data">
                <h4>Departments</h4>
                <p class="stat-number">{{ totalDepartments }}</p>
                <span class="trend font-gray">No change</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="card-icon orange-bg">👤</div>
              <div class="card-data">
                <h4>Present Today</h4>
                <p class="stat-number">{{ presentToday }}</p>
                <span class="trend font-orange">● 81.67%</span>
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
                <canvas id="attendanceChart" ref="attendanceChartRef"></canvas>
              </div>
            </div>

            <div class="grid-card">
              <h3>Recent Leave Requests</h3>
              <div class="leave-list">
                <div class="leave-item">
                  <div class="leave-user">Kavindu Perera</div>
                  <span class="leave-badge pending">Pending</span>
                </div>
                <div class="leave-item">
                  <div class="leave-user">Nimesha Fernando</div>
                  <span class="leave-badge pending">Pending</span>
                </div>
                <div class="leave-item">
                  <div class="leave-user">Tharushi Silva</div>
                  <span class="leave-badge approved">Approved</span>
                </div>
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
                  <tr>
                    <td>EMP001</td>
                    <td>Dhanushka Wickramasinghe</td>
                    <td>IT</td>
                    <td><span class="status-t active">Active</span></td>
                  </tr>
                  <tr>
                    <td>EMP002</td>
                    <td>Hasini Fernando</td>
                    <td>HR</td>
                    <td><span class="status-t active">Active</span></td>
                  </tr>
                  <tr>
                    <td>EMP003</td>
                    <td>Kasun Perera</td>
                    <td>Finance</td>
                    <td><span class="status-t active">Active</span></td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="grid-card">
              <h3>Department Distribution</h3>
              <div class="chart-container pie-size">
                <canvas id="deptPieChart" ref="deptPieChartRef"></canvas>
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
import { apiGet } from '@/api';

const route = useRoute();
const router = useRouter();
const attendanceChartRef = ref(null);
const deptPieChartRef = ref(null);
const totalEmployees = ref(0);
const totalDepartments = ref(0);
const presentToday = ref(0);
const pendingLeaves = ref(0);
const dashboardError = ref('');
const authUser = ref({});
const displayName = computed(() => authUser.value.full_name || 'Admin');
const userPosition = computed(() => authUser.value.position || 'HR Manager');
const chartInstances = {
  attendance: null,
  deptPie: null
};

const loadDashboardSummary = async () => {
  try {
    const data = await apiGet('/api/dashboard-summary');
    totalEmployees.value = data.total_employees ?? 0;
    totalDepartments.value = data.total_departments ?? 0;
    presentToday.value = data.present_today ?? 0;
    pendingLeaves.value = data.pending_leaves ?? 0;
  } catch (error) {
    console.error('Unable to load dashboard summary', error);
    dashboardError.value = error.message;
  }
};

const logout = () => {
  localStorage.removeItem('authToken');
  localStorage.removeItem('authUser');
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

  const ctxLine = attendanceChartRef.value;
  if (ctxLine) {
    chartInstances.attendance = new Chart(ctxLine, {
      type: 'line',
      data: {
        labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
        datasets: [{
          label: 'Attendance',
          data: [85, 98, 92, 90, 96, 35, 70],
          borderColor: '#10b981',
          backgroundColor: 'rgba(16, 185, 129, 0.05)',
          tension: 0.4,
          fill: true,
          borderWidth: 3,
          pointBackgroundColor: '#10b981'
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
  if (ctxPie) {
    chartInstances.deptPie = new Chart(ctxPie, {
      type: 'doughnut',
      data: {
        labels: ['IT', 'HR', 'Finance', 'Marketing', 'Other'],
        datasets: [{
          data: [35, 20, 15, 15, 15],
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
  const savedUser = localStorage.getItem('authUser');
  if (savedUser) {
    try {
      authUser.value = JSON.parse(savedUser);
    } catch (error) {
      console.error('Unable to read signed-in user', error);
    }
  }
  if (route.path === '/dashboard') {
    await loadDashboardSummary();
    createCharts();
  }
});

watch(
  () => route.fullPath,
  (newPath) => {
    if (newPath === '/dashboard') {
      dashboardError.value = '';
      loadDashboardSummary();
      createCharts();
    } else {
      destroyCharts();
    }
  }
);

onBeforeUnmount(() => {
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