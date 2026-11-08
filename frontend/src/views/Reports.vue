<template>
  <div class="reports-container">
    <p v-if="loadError" role="alert" class="api-error">{{ loadError }}</p>
    <EmptyState v-if="loading" message="Loading reports..." />
    <div class="page-header">
      <div class="header-title">
        <h1>📊 Advanced Reports & Analytics</h1>
        <p class="breadcrumb">Dashboard > Reports</p>
      </div>
      <div class="header-actions">
        <button class="btn-primary" @click="exportReports">📥 Export All Reports</button>
      </div>
    </div>

    <div class="analytics-grid">
      <div class="analytics-card">
        <h3>Monthly Attendance Rate</h3>
        <div class="chart-dummy">
          <div v-if="hasAttendanceRate" class="progress-circle" :style="{ background: `conic-gradient(#10b981 ${attendanceRate}%, #e2e8f0 0)` }">
            <span class="circle-value">{{ formatPercent(summary.attendance_rate) }}</span>
          </div>
          <EmptyState v-else :message="loading ? 'Loading attendance data...' : 'No attendance data yet'" height="130px" />
        </div>
        <p v-if="hasAttendanceRate" class="card-desc">Attendance rate from the latest report data.</p>
      </div>

      <div class="analytics-card">
        <h3>Department Wise Performance (%)</h3>
        <div class="progress-list">
          <div class="progress-item" v-for="item in summary.performance" :key="item.division">
            <div class="item-info"><span>{{ item.division }}</span><span>{{ formatPercent(item.value) }}</span></div>
            <div class="bar-bg"><div class="bar-fill" :style="{ width: `${clampPercent(item.value)}%` }"></div></div>
          </div>
          <EmptyState v-if="loading || !summary.performance.length" :message="loading ? 'Loading performance data...' : 'No performance data yet'" />
        </div>
      </div>
    </div>

    <div class="report-table-section">
      <h3>Available Reports for Download</h3>
      <div class="report-table-scroll">
        <table class="report-table">
          <thead>
            <tr>
              <th>Report Name</th>
              <th>Type</th>
              <th>Generated Date</th>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="report in reports" :key="report.name">
              <td><strong>{{ report.name }}</strong></td>
              <td><span class="type-badge" :class="report.type.toLowerCase() === 'pdf' ? 'pdf' : 'excel'">{{ report.type }}</span></td>
              <td>{{ report.generated_date }}</td>
              <td><span class="status-ready">{{ report.status }}</span></td>
              <td><button class="btn-download" @click="downloadReport(report)">Download</button></td>
            </tr>
            <tr v-if="!loading && !reports.length">
              <td colspan="5"><EmptyState message="No reports yet" /></td>
            </tr>
            <tr v-if="loading">
              <td colspan="5"><EmptyState message="Loading reports..." /></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { apiGet } from '@/api';
import EmptyState from '@/components/EmptyState.vue';

const summary = ref({ attendance_rate: null, performance: [] });
const reports = ref([]);
const loadError = ref('');
const loading = ref(false);
const attendanceRate = computed(() => {
  const value = Number(summary.value.attendance_rate);
  return Number.isFinite(value) ? Math.min(100, Math.max(0, value)) : 0;
});
const hasAttendanceRate = computed(() => summary.value.attendance_rate !== null && summary.value.attendance_rate !== undefined && Number.isFinite(Number(summary.value.attendance_rate)));
const clampPercent = value => {
  const number = Number(value ?? 0);
  return Number.isFinite(number) ? Math.min(100, Math.max(0, number)) : 0;
};
const formatPercent = value => `${clampPercent(value).toFixed(1)}%`;

const downloadFile = (name, contents, type = 'text/csv;charset=utf-8') => {
  const url = URL.createObjectURL(new Blob([contents], { type }));
  const link = document.createElement('a');
  link.href = url;
  link.download = name;
  link.click();
  URL.revokeObjectURL(url);
};

const escapeCsv = value => `"${String(value ?? '').replaceAll('"', '""')}"`;
const exportReports = () => {
  const headers = ['name', 'type', 'generated_date', 'status'];
  const rows = reports.value.map(report => headers.map(key => escapeCsv(report[key])).join(','));
  downloadFile('hr-reports.csv', [headers.join(','), ...rows].join('\r\n'));
};

const createPdf = (report) => {
  const escapePdf = value => String(value ?? '')
    .replace(/[^\x20-\x7e]/g, '?')
    .replace(/[\\()]/g, '\\$&');
  const lines = [
    report.name,
    `Type: ${report.type}`,
    `Generated: ${report.generated_date}`,
    `Status: ${report.status}`,
  ];
  const stream = [
    'BT',
    '/F1 16 Tf',
    '50 780 Td',
    `(${escapePdf(lines[0])}) Tj`,
    '/F1 11 Tf',
    ...lines.slice(1).map(line => `0 -24 Td (${escapePdf(line)}) Tj`),
    'ET',
  ].join('\n');
  const objects = [
    '<< /Type /Catalog /Pages 2 0 R >>',
    '<< /Type /Pages /Kids [3 0 R] /Count 1 >>',
    '<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>',
    '<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>',
    `<< /Length ${stream.length} >>\nstream\n${stream}\nendstream`,
  ];
  let pdf = '%PDF-1.4\n';
  const offsets = [0];
  for (const [index, object] of objects.entries()) {
    offsets.push(pdf.length);
    pdf += `${index + 1} 0 obj\n${object}\nendobj\n`;
  }
  const xrefOffset = pdf.length;
  pdf += `xref\n0 ${objects.length + 1}\n0000000000 65535 f \n`;
  pdf += offsets.slice(1).map(offset => `${String(offset).padStart(10, '0')} 00000 n \n`).join('');
  pdf += `trailer\n<< /Size ${objects.length + 1} /Root 1 0 R >>\nstartxref\n${xrefOffset}\n%%EOF`;
  return pdf;
};

const downloadReport = (report) => {
  const fileName = report.name.replace(/[^a-z0-9-_]+/gi, '-');
  if (report.type.toLowerCase() === 'pdf') {
    downloadFile(`${fileName}.pdf`, createPdf(report), 'application/pdf');
    return;
  }
  const columns = ['name', 'type', 'generated_date', 'status'];
  downloadFile(`${fileName}.csv`, [columns.join(','), columns.map(key => escapeCsv(report[key])).join(',')].join('\r\n'));
};

const loadReports = async () => {
  loading.value = true;
  loadError.value = '';
  try {
    const data = await apiGet('/api/reports');
    summary.value = {
      attendance_rate: data.summary?.attendance_rate ?? null,
      performance: Array.isArray(data.summary?.performance) ? data.summary.performance : [],
    };
    reports.value = Array.isArray(data.reports) ? data.reports : [];
  } catch (error) {
    console.error('Failed to load reports', error);
    loadError.value = error.message;
    summary.value = { attendance_rate: null, performance: [] };
    reports.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(loadReports);
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

.reports-container {
  padding: 28px;
  font-family: 'Inter', sans-serif;
  background-color: #f8fafc;
  min-height: 100%;
  min-width: 0;
  max-width: 100%;
  box-sizing: border-box;
}
.api-error { color: #b91c1c; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 10px 14px; }
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
}
.header-title h1 {
  font-size: 24px;
  font-weight: 700;
  color: #062c22;
  margin: 0;
}
.breadcrumb { font-size: 12px; color: #64748b; margin-top: 4px; }
.btn-primary {
  background-color: #062c22;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 25px;
  font-weight: 500;
  cursor: pointer;
}
.btn-primary:hover { background-color: #0c5240; }
.header-actions { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.header-actions { min-width: 0; max-width: 100%; }
.btn-secondary { background: white; color: #062c22; border: 1px solid #cbd5e1; padding: 9px 16px; border-radius: 25px; font-weight: 600; cursor: pointer; }
.btn-secondary:hover { background: #ecfdf5; border-color: #86c7ac; }
.report-create-form { margin: 0 0 24px; padding: 24px; background: white; border: 1px solid #e2e8f0; border-radius: 15px; box-shadow: 0 4px 12px rgba(15, 23, 42, 0.04); }
.form-title { margin: 0 0 20px; color: #062c22; font-size: 17px; font-weight: 700; }
.form-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(190px, 1fr)); gap: 18px; margin-bottom: 22px; }
.form-grid label { display: flex; flex-direction: column; gap: 8px; color: #334155; font-size: 13px; font-weight: 600; }
.form-grid input, .form-grid select { width: 100%; min-height: 42px; box-sizing: border-box; padding: 10px 12px; border: 1px solid #cbd5e1; border-radius: 8px; background: white; color: #0f172a; font: inherit; font-weight: 400; }
.form-grid input:focus, .form-grid select:focus { outline: 2px solid rgba(16, 185, 129, 0.16); border-color: #0c7a55; }
.form-actions { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.submit-btn, .cancel-btn { min-height: 40px; padding: 9px 17px; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; }
.submit-btn { color: white; background: #062c22; border: 1px solid #062c22; }
.submit-btn:hover { background: #0c5240; border-color: #0c5240; }
.cancel-btn { color: #062c22; background: white; border: 1px solid #b7c8c1; }
.cancel-btn:hover { background: #f0fdf4; }
.analytics-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 22px;
  margin-bottom: 26px;
}
.analytics-card {
  min-width: 0;
  background: white;
  border: 1px solid #e8edf2;
  border-radius: 15px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}
.analytics-card h3 { font-size: 15px; color: #062c22; margin : 0 0 15px; font-weight: 600; }
.chart-dummy { display: flex; justify-content: center; padding: 15px 0; min-height: 160px; }
.progress-circle {
  width: 100px; height: 100px;
  border-radius: 50%;
  background: #e2e8f0;
  display: flex; align-items: center; justify-content: center;
}
.circle-value {
  background: white; width: 84px; height: 84px;
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-weight: 700; color: #062c22; font-size: 18px;
}
.card-desc { font-size: 12px; color: #64748b; text-align: center; margin-top: 15px; }

/* Progress Bars */
.progress-list { display: flex; flex-direction: column; gap: 15px; }
.item-info { display: flex; justify-content: space-between; font-size: 12px; color: #475569; font-weight: 500; margin-bottom: 5px; }
.bar-bg { background: #e2e8f0; height: 8px; border-radius: 4px; width: 100%; }
.bar-fill { background: #10b981; height: 100%; border-radius: 4px; }

/* Table */
.report-table-section {
  min-width: 0;
  max-width: 100%;
  background: white; border: 1px solid #e8edf2; border-radius: 15px; padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}
.report-table-section h3 { font-size: 16px; color: #062c22; margin: 0 0 18px; font-weight: 700; }
.report-table-scroll { width: 100%; max-width: 100%; min-width: 0; overflow-x: auto; }
.report-table { width: 100%; min-width: 600px; border-collapse: collapse; font-size: 13px; text-align: left; }
.report-table th { padding: 14px 16px; background: #f8fafc; color: #475569; font-weight: 600; white-space: nowrap; }
.report-table td { padding: 15px 16px; border-bottom: 1px solid #f1f5f9; color: #334155; }
.report-table tbody tr:last-child td { border-bottom: 0; }
.report-table tbody tr:hover { background: #fbfdfc; }
.type-badge { padding: 3px 8px; border-radius: 4px; font-size: 10px; font-weight: 700; }
.type-badge.pdf, .type-badge.excel { background: #d1fae5; color: #065f46; }
.status-ready { color: #10b981; font-weight: 600; font-size: 12px; }
.btn-download { background: #062c22; color: white; border: none; padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.btn-download:hover { background: #0c5240; }
@media (max-width: 760px) {
  .reports-container { padding: 18px; }
  .page-header { align-items: flex-start; flex-direction: column; }
  .header-actions, .header-actions button { width: 100%; }
  .analytics-grid { grid-template-columns: minmax(0, 1fr); }
}
@media (max-width: 480px) {
  .reports-container { padding: 14px; }
  .report-create-form, .analytics-card, .report-table-section { padding: 18px; }
  .header-actions { width: 100%; }
  .header-actions button { flex: 1; }
}
</style>