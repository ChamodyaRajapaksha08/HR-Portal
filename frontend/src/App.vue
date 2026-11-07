<template>
  <div class="app-theme">
    <router-view />
  </div>
</template>

<style>

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
  background-color: #f4f7f6;
  color: #18352b;
  color-scheme: light;
}

.app-theme {
  min-height: 100vh;
  background: #f4f7f6;
  color: #18352b;
  --theme-bg: #f4f7f6;
  --theme-surface: #ffffff;
  --theme-surface-raised: #f0f7f4;
  --theme-border: #d7e5df;
  --theme-text: #18352b;
  --theme-muted: #587067;
  --theme-forest: #386657;
  --theme-forest-dark: #244b3e;
  --theme-emerald: #10b981;
  --theme-emerald-dark: #047857;
  --theme-emerald-pale: #d1fae5;
}

/* Global overrides keep the shared palette consistent with view-scoped styles. */
.app-theme.app-theme .dashboard-layout {
  background-color: var(--theme-bg) !important;
  color: var(--theme-text) !important;
}

.app-theme.app-theme :is(.main-content, .content-body, .page-container, .reports-container, .payroll-container, .settings-container) {
  background-color: #ffffff !important;
  color: var(--theme-text) !important;
}

.app-theme.app-theme :is(.sidebar, .sidebar-brand) {
  background-color: var(--theme-forest) !important;
  color: #f8fafc !important;
  border-color: rgba(255, 255, 255, 0.14) !important;
}

.app-theme.app-theme .sidebar {
  border-radius: 0 !important;
  box-shadow: 4px 0 18px rgba(24, 53, 43, 0.12) !important;
}

.app-theme.app-theme .menu-item.router-link-active {
  background: rgba(16, 185, 129, 0.18) !important;
  color: #d1fae5 !important;
  border-left-color: #34d399 !important;
}

.app-theme.app-theme .menu-item:hover {
  background: rgba(209, 250, 229, 0.12) !important;
  color: #f8fafc !important;
}

.app-theme.app-theme :is(.stat-card, .grid-card, .emp-stat-card, .metric-card-box,
  .table-main-wrapper, .table-controls-card, .chart-card, .analytics-card, .report-table-section,
  .report-create-form, .table-section, .p-card, .settings-card, .login-window) {
  background-color: #ffffff !important;
  color: var(--theme-text) !important;
  border-color: var(--theme-border) !important;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08) !important;
}

.app-theme.app-theme .top-header {
  height: 64px;
  min-height: 64px;
  padding: 0 30px;
  border-bottom: 1px solid var(--theme-border) !important;
  box-shadow: 0 1px 3px rgba(24, 53, 43, 0.08) !important;
}

.app-theme.app-theme :is(.create-record-form, .create-employee-form, .upcoming-item, .leave-item,
  .table-header, .date-picker-mock, .search-input, .search-input-element,
  .search-bar, .form-input, .form-select, .form-grid input, .form-grid select) {
  background-color: #ffffff !important;
  color: var(--theme-text) !important;
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme .bar-bg {
  background-color: #dcebe4 !important;
}

.app-theme.app-theme :is(h1, h2, h3, h4, h5, p, label, th, td, strong,
  .stat-number, .stat-num, .metric-num-val, .card-value, .emp-name, .user-fullname,
  .dept-name-td, .leave-user, .u-name, .item-info, .progress-info-text,
  .card-info h3, .card-data h4, .stat-info h4, .metric-info-data h4,
  .header-title h1, .top-header h2, .settings-card h3, .report-table-section h3) {
  color: var(--theme-text) !important;
}

.app-theme.app-theme :is(.breadcrumb, .subtitle, .user-role, .user-email, .emp-email,
  .card-desc, .section-desc, .pagination-text-info, .pagination-text-info,
  .date-td, .head-name-td, .emp-count-td, .no-upcoming-leaves, .form-grid label,
  .search-icon, .search-lens, .search-bar::placeholder, input::placeholder,
  .form-subtitle, .contact-admin-text, .copyright-text, .input-tip) {
  color: var(--theme-muted) !important;
}

.app-theme.app-theme .sidebar :is(h1, h2, h3, h4, h5, p, span, a, button) {
  color: #f8fafc !important;
}

.app-theme.app-theme :is(.top-header, .sidebar-footer, .settings-card h3, .table-pagination-bar) {
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme table {
  background: #ffffff !important;
  color: var(--theme-text) !important;
}

.app-theme.app-theme table :is(th, td) {
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme table th {
  background: #edf5f1 !important;
  color: var(--theme-text) !important;
}

.app-theme.app-theme table tbody tr:hover {
  background: #f3f8f5 !important;
}

.app-theme.app-theme :is(input:not([type="checkbox"]):not([type="radio"]), select, textarea, .dropdown-filter, .filter-dropdown) {
  background-color: #ffffff !important;
  color: var(--theme-text) !important;
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme :is(input:not([type="checkbox"]):not([type="radio"]), textarea)::placeholder {
  color: var(--theme-muted) !important;
  opacity: 1;
}

.app-theme.app-theme :is(input, select, textarea):focus {
  border-color: var(--theme-emerald) !important;
  outline-color: rgba(16, 185, 129, 0.28) !important;
}

.app-theme.app-theme :is(.btn-primary, .submit-btn, .btn-save, .add-emp-btn,
  .add-record-btn, .add-department-btn, .btn-download, .btn-payslip, .filter-btn,
  .login-submit-btn) {
  background-color: var(--theme-emerald-dark) !important;
  border-color: var(--theme-emerald-dark) !important;
  color: #f8fafc !important;
}

.app-theme.app-theme :is(.btn-primary, .submit-btn, .btn-save, .add-emp-btn,
  .add-record-btn, .add-department-btn, .btn-download, .btn-payslip, .filter-btn,
  .login-submit-btn):hover {
  background-color: #059669 !important;
  border-color: #059669 !important;
}

.app-theme.app-theme :is(.btn-secondary, .cancel-btn, .clear-filter-btn, .table-icon-btn,
  .nav-arrow-btn, .nav-number-btn, .social-btn) {
  background-color: #f0f7f4 !important;
  color: var(--theme-text) !important;
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme :is(.btn-secondary, .cancel-btn, .clear-filter-btn, .table-icon-btn,
  .nav-arrow-btn, .nav-number-btn, .social-btn):hover {
  background-color: #dcebe4 !important;
}

.app-theme.app-theme :is(.blue-bg, .light-green-bg, .blue-tag, .it-tag, .dept-badge,
  .status-t.active, .badge-active, .type-badge.excel, .status-badge.approved,
  .status-badge.paid, .status-badge-pill.active-bg, .status-badge-pill.present-bg,
  .status-pill.active, .days-badge.b-blue, .leave-badge.approved) {
  background-color: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.green-bg, .hr-tag, .status-done, .positive, .status-ready) {
  background-color: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.orange-bg, .orange-tag, .status-pending, .status-badge.pending,
  .status-badge-pill.late-bg) {
  background-color: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.red-bg, .status-badge.rejected, .status-badge-pill.absent-bg,
  .type-badge.pdf, .pink-bg, .pink-tag) {
  background-color: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.purple-bg, .days-badge.b-purple) {
  background-color: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.status-ready, .text-green, .text-blue, .trend.font-green,
  .trend.font-pink, .font-pink, .text-red, .text-orange, .font-orange,
  .trend.font-orange, .edit-color, .delete-color, .positive, .status-done,
  .pagination-info) {
  color: var(--theme-emerald-dark) !important;
}

.app-theme.app-theme :is(.role-text, .id-code, .emp-id) {
  color: #334155 !important;
}

.app-theme.app-theme .api-error {
  background: #fef2f2 !important;
  border-color: #fecaca !important;
  color: #991b1b !important;
}

.app-theme.app-theme :is(.status-badge.inactive-bg, .status-badge-pill.inactive-bg,
  .status-pill.on-leave, .status-badge.pending, .leave-badge.pending,
  .status-badge-pill.late-bg, .status-pending) {
  background: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.status-badge.inactive, .status-badge-pill.inactive-bg,
  .status-pill.inactive, .status-t.inactive) {
  background: #e5f3ec !important;
  color: #386657 !important;
}

.app-theme.app-theme :is(.finance-tag, .orange-tag) {
  background: #d1fae5 !important;
  color: #065f46 !important;
}

.app-theme.app-theme :is(.page-num-btn.active, .nav-number-btn.active-page,
  .page-number.active, .pagination-button.active) {
  background: var(--theme-emerald-dark) !important;
  color: #fff !important;
  border-color: var(--theme-emerald-dark) !important;
}

.app-theme.app-theme .login-container .input-field {
  background: #ffffff !important;
  color: #0f172a !important;
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme :is(.item-info) {
  color: var(--theme-muted) !important;
}

.app-theme.app-theme :is(.contact-link, .forgot-link) {
  color: var(--theme-emerald-dark) !important;
}

.app-theme.app-theme .login-container {
  background-color: var(--theme-bg) !important;
}

.app-theme.app-theme .form-panel {
  background-color: #f0faf5 !important;
}

.app-theme.app-theme .login-window {
  background: #ffffff !important;
  border-color: var(--theme-border) !important;
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.18) !important;
}

.app-theme.app-theme :is(.form-center-box h2, .login-container label, .remember-me) {
  color: var(--theme-text) !important;
}

.app-theme.app-theme .login-container .form-divider::before {
  background: var(--theme-border) !important;
}

.app-theme.app-theme .form-divider span {
  background: #f0faf5 !important;
  color: var(--theme-muted) !important;
}

.app-theme.app-theme :is(.emp-avatar, .user-avatar-small) {
  background: #dcebe4 !important;
  color: #244b3e !important;
}

.app-theme.app-theme :is(.table-responsive, .table-scroll-area, .report-table-scroll) {
  border-color: var(--theme-border) !important;
}

.app-theme.app-theme ::selection {
  background: rgba(16, 185, 129, 0.3);
  color: #fff;
}

.app-theme.app-theme * {
  scrollbar-color: #8eaa9c #e8f1ec;
  scrollbar-width: thin;
}

@media (max-width: 700px) {
  .app-theme.app-theme .sidebar {
    width: 210px;
  }
}
</style>