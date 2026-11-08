import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'

import Employees from '../views/Employees.vue'
import Departments from '../views/Departments.vue'
import Attendance from '../views/Attendance.vue'
import Leaves from '../views/Leaves.vue'

const managementRoles = ['Admin', 'HR Manager']
const employeeRoutes = [
  { path: 'profile', name: 'MyProfile', component: () => import('../views/MyProfile.vue'), meta: { roles: ['Employee'] } },
  { path: 'my-attendance', name: 'MyAttendance', component: () => import('../views/MyAttendance.vue'), meta: { roles: ['Employee'] } },
  { path: 'my-leaves', name: 'MyLeaves', component: () => import('../views/MyLeaves.vue'), meta: { roles: ['Employee'] } },
  { path: 'my-payslips', name: 'MyPayslips', component: () => import('../views/MyPayslips.vue'), meta: { roles: ['Employee'] } },
  { path: 'change-password', name: 'ChangePassword', component: () => import('../views/ChangePassword.vue'), meta: { roles: ['Employee'] } },
]

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  {
    path: '/dashboard',
    component: Dashboard,
    meta: { requiresAuth: true },
    children: [
      { path: 'employees', component: Employees, meta: { roles: managementRoles } },
      { path: 'departments', component: Departments, meta: { roles: managementRoles } },
      { path: 'attendance', component: Attendance, meta: { roles: managementRoles } },
      { path: 'leaves', component: Leaves, meta: { roles: managementRoles } },
      {
        path: 'payroll',
        name: 'Payroll',
        component: () => import('../views/Payroll.vue'),
        meta: { roles: managementRoles }
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('../views/Reports.vue'),
        meta: { roles: managementRoles }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue'),
        meta: { roles: managementRoles }
      },
      ...employeeRoutes,
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const token = localStorage.getItem('authToken')
  const isAuthenticated = Boolean(token)
  let user = {}
  try {
    user = JSON.parse(localStorage.getItem('authUser') || '{}')
  } catch {
    user = {}
  }
  const role = user.role
  const home = '/dashboard'

  if (to.matched.some(record => record.meta.requiresAuth) && !isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.path === '/login' && isAuthenticated) {
    return home
  }

  const allowedRoles = to.matched.map(record => record.meta.roles).filter(Boolean).flat()
  if (allowedRoles.length && !allowedRoles.includes(role)) {
    return home
  }
})

export default router