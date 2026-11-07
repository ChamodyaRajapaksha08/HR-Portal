import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'

import Employees from '../views/Employees.vue'
import Departments from '../views/Departments.vue'
import Attendance from '../views/Attendance.vue'
import Leaves from '../views/Leaves.vue'

const DashboardOverview = {
  template: `
    <div class="stats-grid">
      <div class="stat-card"><h4>Total Employees</h4><p class="stat-number">124</p></div>
      <div class="stat-card"><h4>On Leave Today</h4><p class="stat-number alert">5</p></div>
      <div class="stat-card"><h4>Total Departments</h4><p class="stat-number">6</p></div>
    </div>
  `
}

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  {
    path: '/dashboard',
    component: Dashboard,
    children: [
      { path: '', component: DashboardOverview },
      { path: 'employees', component: Employees },
      { path: 'departments', component: Departments },
      { path: 'attendance', component: Attendance },
      { path: 'leaves', component: Leaves },
      {
        path: 'payroll',
        name: 'Payroll',
        component: () => import('../views/Payroll.vue')
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('../views/Reports.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const isAuthenticated = Boolean(localStorage.getItem('authToken'))
  const requiresAuthentication = to.matched.some((record) => record.path.startsWith('/dashboard'))

  if (requiresAuthentication && !isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.path === '/login' && isAuthenticated) {
    return '/dashboard'
  }
})

export default router