<script>
  import { scanHistory } from '../stores';
  import HistoryDetailModal from '../components/historyPage/HistoryDetailModal.svelte';
  import ConfirmationModal from '../components/general/ConfirmationModal.svelte';
  
  let selectedScan = null;
  let showDetailModal = false;

  let showConfirmModal = false;
  
  function viewScanDetail(scan) {
    selectedScan = scan;
    showDetailModal = true;
  }
  
  function closeDetailModal() {
    showDetailModal = false;
  }

  let titleConfirmModal
  let confirmMessage
  function viewDeleteConfirmationScanData(scan) {
    showConfirmModal = true;
    selectedScan = scan;

    titleConfirmModal = `Delete Scan ${scan.name}`;
    confirmMessage = `Are you sure you want to delete the scan with title ${scan.name}? This action cannot be undone.`;
  }

  function deleteScanData() {
    // Remove from history (this should ideally be done after successful deletion)
    $scanHistory = $scanHistory.filter(s => s.id !== selectedScan.id);
    
    closeConfirmModal();
  }

  function closeConfirmModal() {
    showConfirmModal = false;
  }
  
  function formatDate(dateString) {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric'
    }).format(date);
  }
  
  // Define status badge class based on status
  function getBadgeClass(status) {
    switch(status) {
      case 'completed':
        return 'badge-success';
      case 'failed':
        return 'badge-error';
      case 'running':
        return 'badge-info';
      default:
        return 'badge-warning';
    }
  }
</script>

<div class="p-6">
  <h1 class="text-3xl font-bold mb-8 text-primary">Scan History</h1>
  
  <!-- Search and filter (placeholder for future expansion) -->
  <!-- <div class="mb-6">
    <div class="join">
      <input type="text" placeholder="Search scans..." class="join-item input input-bordered w-full max-w-xs" />
      <button class="join-item btn btn-primary">Search</button>
    </div>
  </div> -->
  
  <!-- History Table -->
  <div class="overflow-x-auto">    <table class="table table-zebra w-full">
      <thead>
        <tr class="text-base-content bg-base-300">
          <th>ID</th>
          <th>Project Name</th>
          <th>Date</th>
          <th>Status</th>
          <th>Vulnerabilities</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if $scanHistory.length === 0}
          <tr>
            <td colspan="6" class="text-center text-black text-bold py-8">No scan history found. Try running a scan first.</td>
          </tr>
        {:else}
          {#each $scanHistory as scan (scan.id)}
            <tr>
              <td class="font-medium text-base-content">{$scanHistory.indexOf(scan) + 1}</td>
              <td class="font-medium text-base-content">{scan.name}</td>
              <td class="text-base-content">{formatDate(scan.date)}</td>
              <td>
                <span class="badge {getBadgeClass(scan.status)} font-medium">{scan.status}</span>
              </td>
              <td>
                {#if scan.vulnerabilities > 0}
                  <span class="text-error font-bold">{scan.vulnerabilities}</span>
                {:else}
                  <span class="text-success font-bold">{scan.vulnerabilities}</span>
                {/if}
              </td>
              <td>
                <button class="btn btn-sm btn-primary" on:click={() => viewScanDetail(scan)}>
                  Details
                </button>
                <button class="btn btn-sm btn-error" on:click={() => viewDeleteConfirmationScanData(scan)}>
                  Delete
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  
  <!-- History Detail Modal -->
  {#if showDetailModal && selectedScan}
    <HistoryDetailModal scan={selectedScan} on:close={closeDetailModal} />
  {/if}

  <!-- Confirmation Modal -->
  {#if showConfirmModal && selectedScan}
    <ConfirmationModal 
    on:confirm={deleteScanData} 
    on:cancel={closeConfirmModal} 
    TitleConfirmationModal={titleConfirmModal} 
    MessageConfirmationModal={confirmMessage} 
    />
  {/if}
</div>
