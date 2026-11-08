<template>
  <section class="employee-home">
    <p v-if="error" class="api-error" role="alert">{{ error }}</p>
    <EmptyState v-if="loading" message="Loading your dashboard..." />

    <template v-else>
    <div class="summary-grid">
      <article v-for="balance in balances" :key="balance.leave_type" class="summary-card">
        <span>{{ balance.leave_type }} Leave</span>
        <strong>{{ balance.remaining }} <small>remaining</small></strong>
        <p>{{ balance.used }} used of {{ balance.entitlement }} days</p>
      </article>
    </div>

    <div class="summary-grid attendance-summary">
      <article v-for="item in attendanceStats" :key="item.label" class="summary-card">
        <span>This Month · {{ item.label }}</span>
        <strong>{{ item.value }}</strong>
      </article>
    </div>

    <div class="home-grid">
      <section class="home-card">
        <h3>Latest Leave Requests</h3>
        <div v-for="leave in latestLeaves" :key="leave.id" class="list-row">
          <div>
            <strong>{{ leave.leave_type }}</strong>
            <small>{{ leave.from_date }} – {{ leave.to_date }}</small>
          </div>
          <span class="status-badge" :class="statusClass(leave.status)">{{ leave.status }}</span>
        </div>
        <EmptyState v-if="!loading && !latestLeaves.length" message="No leave requests yet" />
      </section>

      <section class="home-card">
        <h3>Latest Payslip</h3>
        <template v-if="latestPayslip">
          <p class="payslip-amount">{{ money(latestPayslip.net_salary) }}</p>
          <p class="payslip-meta">{{ payslipMonth(latestPayslip) }} · {{ latestPayslip.status }}</p>
          <router-link to="/dashboard/my-payslips" class="view-link">View payslips</router-link>
        </template>
        <EmptyState v-else-if="!loading" message="No payslips yet" />
      </section>
    </div>
    </template>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { apiGet } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const loading = ref(true);
const error = ref('');
const balances = ref([]);
const attendanceSummary = ref({});
const leaves = ref([]);
const payslips = ref([]);
const latestLeaves = computed(() => leaves.value.slice(0, 5));
const latestPayslip = computed(() => payslips.value[0] || null);
const attendanceStats = computed(() => [
  { label: 'Present', value: Number(attendanceSummary.value.present) || 0 },
  { label: 'Late', value: Number(attendanceSummary.value.late) || 0 },
  { label: 'Absent', value: Number(attendanceSummary.value.absent) || 0 },
]);

const money = value => `LKR ${(Number(value) || 0).toLocaleString()}`;
const statusClass = status => (status || '').toLowerCase();
const payslipMonth = payslip => payslip.paid_date
  ? new Date(`${payslip.paid_date}T00:00:00`).toLocaleDateString(undefined, { month: 'long', year: 'numeric' })
  : 'Unpaid';

onMounted(async () => {
  try {
    const [balanceData, attendanceData, leaveData, payslipData] = await Promise.all([
      apiGet('/api/me/leave-balance'),
      apiGet('/api/me/attendance'),
      apiGet('/api/me/leaves'),
      apiGet('/api/me/payslips'),
    ]);
    balances.value = Array.isArray(balanceData.balances) ? balanceData.balances : [];
    attendanceSummary.value = attendanceData.summary || {};
    leaves.value = Array.isArray(leaveData.requests) ? leaveData.requests : [];
    payslips.value = Array.isArray(payslipData.payslips) ? payslipData.payslips : [];
  } catch (requestError) {
    error.value = requestError.message || 'Unable to load your dashboard.';
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.employee-home { display: grid; gap: 20px; }
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; margin: 0; }
.summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 16px; }
.summary-card, .home-card { min-width: 0; padding: 18px; border: 1px solid #e2e8f0; border-radius: 12px; background: #fff; box-shadow: 0 4px 12px rgba(15, 23, 42, .03); }
.summary-card > span { display: block; color: #64748b; font-size: 12px; font-weight: 600; text-transform: uppercase; }
.summary-card strong { display: block; margin-top: 8px; color: #18352b; font-size: 27px; }
.summary-card strong small { font-size: 13px; font-weight: 500; }
.summary-card p { margin: 6px 0 0; color: #64748b; font-size: 13px; }
.home-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; }
.home-card h3 { margin: 0 0 14px; color: #18352b; font-size: 16px; }
.list-row { display: flex; justify-content: space-between; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px solid #f1f5f9; }
.list-row > div { display: grid; gap: 4px; min-width: 0; }
.list-row strong { color: #334155; font-size: 14px; }
.list-row small, .payslip-meta { color: #64748b; font-size: 12px; }
.status-badge { flex: 0 0 auto; padding: 4px 8px; border-radius: 6px; background: #f1f5f9; color: #475569; font-size: 11px; font-weight: 600; }
.status-badge.approved { background: #dcfce7; color: #15803d; }
.status-badge.pending { background: #fef3c7; color: #b45309; }
.status-badge.rejected { background: #fee2e2; color: #b91c1c; }
.payslip-amount { margin: 0 0 6px; color: #18352b; font-size: 30px; font-weight: 700; }
.view-link { display: inline-block; margin-top: 12px; color: #047857; font-size: 13px; font-weight: 600; text-decoration: none; }
@media (max-width: 700px) { .home-grid { grid-template-columns: minmax(0, 1fr); } }
</style>
