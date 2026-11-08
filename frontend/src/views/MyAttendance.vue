<template>
  <section class="self-service-page">
    <header class="page-heading"><div><h2>My Attendance</h2><p class="breadcrumb">Dashboard > My Attendance</p></div></header>
    <p v-if="error" class="api-error" role="alert">{{ error }}</p>
    <EmptyState v-if="loading" message="Loading your attendance..." />
    <template v-else>
      <div class="summary-grid">
        <article v-for="item in summaryCards" :key="item.label" class="summary-card">
          <span>This Month · {{ item.label }}</span><strong>{{ item.value }}</strong>
        </article>
      </div>
      <section class="table-card">
        <div class="table-scroll">
          <table>
            <thead><tr><th>Date</th><th>Department</th><th>Check In</th><th>Check Out</th><th>Status</th></tr></thead>
            <tbody>
              <tr v-for="record in records" :key="record.id">
                <td>{{ formatDate(record.created_at) }}</td><td>{{ record.department || '—' }}</td>
                <td>{{ record.check_in || '—' }}</td><td>{{ record.check_out || '—' }}</td>
                <td><span class="status-badge" :class="(record.status || '').toLowerCase()">{{ record.status || '—' }}</span></td>
              </tr>
              <tr v-if="!records.length"><td colspan="5"><EmptyState message="No attendance records yet" /></td></tr>
            </tbody>
          </table>
        </div>
      </section>
    </template>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { apiGet } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const loading = ref(true);
const error = ref('');
const records = ref([]);
const summary = ref({});
const summaryCards = computed(() => [
  { label: 'Present', value: Number(summary.value.present) || 0 },
  { label: 'Late', value: Number(summary.value.late) || 0 },
  { label: 'Absent', value: Number(summary.value.absent) || 0 },
]);
const formatDate = value => value ? new Date(value).toLocaleDateString() : '—';

onMounted(async () => {
  try {
    const data = await apiGet('/api/me/attendance');
    records.value = Array.isArray(data.attendance) ? data.attendance : [];
    summary.value = data.summary || {};
  } catch (requestError) {
    error.value = requestError.message || 'Unable to load your attendance.';
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.self-service-page { width: 100%; min-width: 0; }
.page-heading { margin-bottom: 20px; }
.page-heading h2 { margin: 0 0 4px; color: #0f172a; font-size: 24px; }
.breadcrumb { margin: 0; color: #64748b; font-size: 13px; }
.api-error { color: #b91c1c; }
.summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(170px, 1fr)); gap: 16px; margin-bottom: 20px; }
.summary-card, .table-card { padding: 18px; border: 1px solid #e2e8f0; border-radius: 12px; background: white; }
.summary-card span { display: block; color: #64748b; font-size: 12px; font-weight: 600; text-transform: uppercase; }
.summary-card strong { display: block; margin-top: 8px; color: #18352b; font-size: 27px; }
.table-scroll { width: 100%; max-width: 100%; overflow-x: auto; }
table { width: 100%; border-collapse: collapse; text-align: left; }
th, td { padding: 12px 14px; border-bottom: 1px solid #e2e8f0; white-space: nowrap; font-size: 13px; }
th { color: #64748b; font-size: 11px; text-transform: uppercase; }
.status-badge { padding: 4px 8px; border-radius: 6px; background: #f1f5f9; color: #475569; font-size: 11px; font-weight: 600; }
.status-badge.present { background: #dcfce7; color: #15803d; }
.status-badge.late { background: #fef3c7; color: #b45309; }
.status-badge.absent { background: #fee2e2; color: #b91c1c; }
</style>
