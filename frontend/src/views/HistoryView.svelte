<script>
  import { onMount } from 'svelte';
  import { scanHistory } from '../stores';
  import HistoryDetailModal from '../components/historyPage/HistoryDetailModal.svelte';
  import ComparisonModal from '../components/historyPage/ComparisonModal.svelte';
  import ConfirmationModal from '../components/general/ConfirmationModal.svelte';
  import InformationModal from '../components/general/InformationModal.svelte';
  
  import { GetScanDatas, DeleteScan, GetRepoDatas, GetScanDetail } from '../../wailsjs/go/main/App';
  
  // State variables
  let scans = [];
  let isLoading = false;
  let currentPage = 1;
  let perPage = 10;
  let totalItems = 0;
  let totalPages = 0;
  let searchTerm = '';
  let sortDesc = false;
  
  // Repository filter
  let selectedRepository = null;
  let repositoryFilter = '';
  let repositories = [];
  let showRepositoryDropdown = false;
  let searchTimeout = null;
  
  // Modal states
  let selectedScan = null;
  let showDetailModal = false;
  let showConfirmModal = false;
  let showInfoModal = false;
  let showComparisonModal = false;
  let infoModalTitle = '';
  let infoModalMessage = '';
  let titleConfirmModal = '';
  let confirmMessage = '';
  let repoId = 0;
  
  // Compare mode states
  let isCompareMode = false;
  let selectedScansForComparison = [];
  let comparisonData = null;
  
  // Computed
  $: totalPages = Math.ceil(totalItems / perPage);
  $: startItem = (currentPage - 1) * perPage + 1;
  $: endItem = Math.min(currentPage * perPage, totalItems);
  
  // Reactive statement to reload data when repoId changes
  $: if (repoId !== undefined) {
    loadScans();
  }
  
  onMount(() => {
    // Only load repositories for filter on mount, reactive statement will handle scans
    loadRepositoriesForFilter();
  });
  
  async function loadScans() {
    isLoading = true;
    try {
      // Use repoId directly for backend filtering - 0 means all repositories
      const response = await GetScanDatas(repoId, currentPage, perPage, searchTerm, sortDesc);
      
      // Handle different possible response formats
      if (Array.isArray(response.data)) {
        scans = response.data;
        totalItems = response.total || response.data.length;
        
        // If we get less data than expected, adjust total
        if (response.data.length < perPage && currentPage === 1) {
          totalItems = response.data.length;
        } else if (response.data.length < perPage) {
          totalItems = (currentPage - 1) * perPage + response.data.length;
        }
      } else if (response && typeof response === 'object') {
        // Fallback for different response structure
        scans = response['data'] || [];
        totalItems = response['total'] || 0;
      } else {
        scans = [];
        totalItems = 0;
      }
      
    } catch (error) {
      showInfo('Error', `Failed to load scan history: ${error.message || error}`);
      scans = [];
      totalItems = 0;
    } finally {
      isLoading = false;
    }
  }
  
  async function loadRepositoriesForFilter() {
    try {
      const userId = 1;
      const response = await GetRepoDatas(userId, 1, 10, '', false);
      repositories = Array.isArray(response) ? response : (response?.data || []);
    } catch (error) {
      repositories = [];
    }
  }
  
  // Repository filter functions
  async function searchRepositories() {
    if (!repositoryFilter.trim()) {
      repositories = [];
      showRepositoryDropdown = false;
      return;
    }
    
    try {
      const userId = 1;
      const response = await GetRepoDatas(userId, 1, 10, repositoryFilter, false);
      repositories = Array.isArray(response) ? response : (response?.data || []);
      showRepositoryDropdown = true;
    } catch (error) {
      repositories = [];
      showRepositoryDropdown = true;
    }
  }
  
  function handleRepositorySearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      searchRepositories();
    }, 300);
  }
  
  function selectRepository(repo) {
    selectedRepository = repo;
    repositoryFilter = repo.name;
    repoId = repo.id || 0; // Set repoId for filtering scans - reactive statement will trigger loadScans
    repositories = [];
    showRepositoryDropdown = false;
    currentPage = 1;
    // Force immediate re-render by clearing scans first
    scans = [];
  }
  
  function clearRepositoryFilter() {
    selectedRepository = null;
    repositoryFilter = '';
    repoId = 0; // Reset repoId to 0 for all repositories - reactive statement will trigger loadScans
    repositories = [];
    showRepositoryDropdown = false;
    currentPage = 1;
    // Force immediate re-render by clearing scans first
    scans = [];
    loadRepositoriesForFilter();
  }
  
  // Pagination functions
  function goToPage(page) {
    if (page >= 1 && page <= totalPages && page !== currentPage) {
      currentPage = page;
      loadScans();
    }
  }
  
  function nextPage() {
    if (currentPage < totalPages) {
      currentPage++;
      loadScans();
    }
  }
  
  function prevPage() {
    if (currentPage > 1) {
      currentPage--;
      loadScans();
    }
  }
  
  // Search function
  function handleSearch() {
    currentPage = 1;
    loadScans();
  }
  
  // Sort function
  function toggleSort() {
    sortDesc = !sortDesc;
    loadScans();
  }
  
  // Compare mode functions
  function toggleCompareMode() {
    isCompareMode = !isCompareMode;
    selectedScansForComparison = [];
  }
  
  function handleScanSelection(scan, event) {
    const isChecked = event.target.checked;
    
    if (isChecked) {
      // Check if we already have 2 scans selected
      if (selectedScansForComparison.length >= 2) {
        event.target.checked = false;
        showInfo('Selection Limit', 'You can only select exactly 2 scans for comparison.');
        return;
      }
      
      // Check if both scans are from the same repository
      if (selectedScansForComparison.length === 1) {
        const firstScan = selectedScansForComparison[0];
        if (firstScan.repository_name !== scan.repository_name) {
          event.target.checked = false;
          showInfo('Repository Mismatch', 'You can only compare scans from the same repository.');
          return;
        }
      }
      
      selectedScansForComparison = [...selectedScansForComparison, scan];
    } else {
      selectedScansForComparison = selectedScansForComparison.filter(s => s.scan.id !== scan.scan.id);
    }
  }
  
  async function startComparison() {
    if (selectedScansForComparison.length !== 2) {
      showInfo('Selection Required', 'Please select exactly 2 scans to compare.');
      return;
    }
    
    try {
      const [scan1, scan2] = selectedScansForComparison;
      
      // Determine which is older and newer based on scan time
      const scan1Time = new Date(scan1.scan.scan_time);
      const scan2Time = new Date(scan2.scan.scan_time);
      
      const olderScan = scan1Time < scan2Time ? scan1 : scan2;
      const newerScan = scan1Time < scan2Time ? scan2 : scan1;
      
      // Fetch detailed data for both scans
      const [olderScanDetails, newerScanDetails] = await Promise.all([
        GetScanDetail(olderScan.scan.result),
        GetScanDetail(newerScan.scan.result)
      ]);
      
      // Process comparison
      comparisonData = processComparison(olderScanDetails, newerScanDetails, olderScan, newerScan);
      showComparisonModal = true;
      
    } catch (error) {
      showInfo('Error', `Failed to load comparison data: ${error.message || error}`);
    }
  }
  
  function processComparison(olderScanJSON, newerScanJSON, olderScanMeta, newerScanMeta) {
    // Extract results arrays from both scans
    const olderResults = olderScanJSON?.results || [];
    const newerResults = newerScanJSON?.results || [];

    // Create fingerprint maps for efficient lookup, with values as arrays (lists)
    const olderFingerprints = new Map();
    const newerFingerprints = new Map();
    
    // Group older results by fingerprint
    olderResults.forEach(result => {
        const fingerprint = result.extra?.fingerprint;
        if (fingerprint) {
            // If the fingerprint isn't in the map, create a new array for it
            if (!olderFingerprints.has(fingerprint)) {
                olderFingerprints.set(fingerprint, []);
            }
            // Push the current result into the array for that fingerprint
            olderFingerprints.get(fingerprint).push(result);
        }
    });
    
    // Group newer results by fingerprint
    newerResults.forEach(result => {
        const fingerprint = result.extra?.fingerprint;
        if (fingerprint) {
            // If the fingerprint isn't in the map, create a new array for it
            if (!newerFingerprints.has(fingerprint)) {
                newerFingerprints.set(fingerprint, []);
            }
            // Push the current result into the array for that fingerprint
            newerFingerprints.get(fingerprint).push(result);
        }
    });

    // Categorize vulnerabilities
    const newVulnerabilities = [];
    const fixedVulnerabilities = [];
    const unresolvedVulnerabilities = [];
    
    // Find new vulnerabilities (in newer but not in older)
    newerFingerprints.forEach((resultsArray, fingerprint) => {
        if (!olderFingerprints.has(fingerprint)) {
            // Since the fingerprint is new, all results in its array are new
            resultsArray.forEach(result => {
                newVulnerabilities.push({
                    ...result,
                    status: 'new',
                    statusColor: 'text-error'
                });
            });
        } else {
            // Since the fingerprint exists in both, all results are unresolved
            resultsArray.forEach(result => {
                unresolvedVulnerabilities.push({
                    ...result,
                    status: 'unresolved',
                    statusColor: 'text-warning'
                });
            });
        }
    });
    
    // Find fixed vulnerabilities (in older but not in newer)
    olderFingerprints.forEach((resultsArray, fingerprint) => {
        if (!newerFingerprints.has(fingerprint)) {
            // Since the fingerprint is gone, all results in its array are fixed
            resultsArray.forEach(result => {
                fixedVulnerabilities.push({
                    ...result,
                    status: 'fixed',
                    statusColor: 'text-success'
                });
            });
        }
    });

    return {
        metadata: {
            olderScan: {
                name: olderScanMeta.repository_name,
                date: olderScanMeta.scan.scan_time,
                totalVulnerabilities: olderResults.length
            },
            newerScan: {
                name: newerScanMeta.repository_name,
                date: newerScanMeta.scan.scan_time,
                totalVulnerabilities: newerResults.length
            }
        },
        summary: {
            new: newVulnerabilities.length,
            fixed: fixedVulnerabilities.length,
            unresolved: unresolvedVulnerabilities.length
        },
        details: {
            new: newVulnerabilities,
            fixed: fixedVulnerabilities,
            unresolved: unresolvedVulnerabilities,
            all: [...newVulnerabilities, ...fixedVulnerabilities, ...unresolvedVulnerabilities]
        }
    };
}
  
  // Modal functions
  async function viewScanDetail(scan) {
    try {
      selectedScan = {
        name: scan.repository_name,
        path: scan.repository_path,
        date: scan.scan.scan_time,
        status: scan.scan.status,
        vulnerabilities: scan.scan.vulnerabilities,
        report: null, // Initialize with null while loading
        isLoading: true
      };

      // Show modal immediately with loading state
      showDetailModal = true;
      
      // Then fetch the details
      const detailReport = await GetScanDetail(scan.scan.result);
      
      // Update the selected scan with the report data
      selectedScan = {
        ...selectedScan,
        report: detailReport,
        isLoading: false
      };
    } catch (error) {

      // Update with error state if already showing modal
      if (showDetailModal && selectedScan) {
        selectedScan = {
          ...selectedScan,
          error: `Failed to load details: ${error.message || 'Unknown error'}`,
          isLoading: false
        };
      } else {
        showInfo('Error', `Failed to load scan details: ${error.message || 'Unknown error'}`);
      }
    }
  }
  
  function closeDetailModal() {
    showDetailModal = false;
  }

  function viewComparisonVulnerabilityDetail(vulnerability) {
    // Create a mock scan object to reuse the existing detail modal
    const mockScan = {
      name: comparisonData?.metadata?.newerScan?.name || 'Comparison View',
      path: vulnerability.path || '',
      date: new Date().toISOString(),
      status: 'completed',
      vulnerabilities: 1,
      report: {
        results: [vulnerability]
      },
      isLoading: false
    };
    
    selectedScan = mockScan;
    showDetailModal = true;
  }

  function viewDeleteConfirmationScanData(scan) {
    showConfirmModal = true;
    selectedScan = scan;
    titleConfirmModal = `Delete Scan "${scan.repo_name || scan.name}"`;
    confirmMessage = `Are you sure you want to delete the scan of "${scan.repo_name || scan.name}"? This action cannot be undone.`;
  }

  async function deleteScanData() {
    if (!selectedScan) return;
    
    try {
      // Use the correct scan ID structure based on the data format
      const scanId = selectedScan.scan?.id || selectedScan.id;
      await DeleteScan(scanId, 1); // scanId first, then userId
      showInfo('Success', 'Scan deleted successfully!');
      loadScans();
    } catch (error) {
      showInfo('Error', `Failed to delete scan: ${error.message || error}`);
    }
    
    closeConfirmModal();
  }

  function closeConfirmModal() {
    showConfirmModal = false;
    selectedScan = null;
  }
  
  function showInfo(title, message) {
    infoModalTitle = title;
    infoModalMessage = message;
    showInfoModal = true;
  }
  
  function formatDate(dateString) {
    if (!dateString) return '-';
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
  
  // Close dropdown when clicking outside
  function handleClickOutside(event) {
    if (!event.target.closest('.repository-filter-container')) {
      showRepositoryDropdown = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="p-6 flex flex-col h-full">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-3xl font-bold text-primary">Scan History</h1>
    
    <!-- Compare Mode Controls -->
    <div class="flex gap-2">
      {#if isCompareMode}
        <button 
          class="btn btn-success"
          disabled={selectedScansForComparison.length !== 2}
          on:click={startComparison}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          Compare Now ({selectedScansForComparison.length}/2)
        </button>
        <button class="btn btn-outline" on:click={toggleCompareMode}>
          Cancel
        </button>
      {:else}
        <button class="btn btn-primary" on:click={toggleCompareMode}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          Compare Results
        </button>
      {/if}
    </div>
  </div>
  
  <!-- Search and Filter Controls -->
  <div class="flex gap-4 mb-6">
    
    <!-- Repository Filter -->
    <div class="repository-filter-container relative">
      <div class="form-control">
        <label class="label" for="repository-filter">
          <span class="label-text font-medium">Filter by Repository</span>
        </label>
        <div class="relative">
          <input 
            id="repository-filter"
            type="text" 
            placeholder={selectedRepository ? selectedRepository.name : "All Repositories"}
            class="input input-bordered w-64 text-black"
            bind:value={repositoryFilter}
            on:input={handleRepositorySearchInput}
            on:focus={() => repositoryFilter && searchRepositories()}
            class:input-success={selectedRepository}
          />
          
          <!-- Clear button -->
          {#if selectedRepository}
            <button 
              type="button"
              class="absolute right-2 top-1/2 transform -translate-y-1/2 btn btn-ghost btn-sm btn-circle"
              on:click={clearRepositoryFilter}
              title="Clear filter"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          {/if}
          
          <!-- Repository Dropdown -->
          {#if showRepositoryDropdown}
            <div class="absolute top-full left-0 right-0 z-10 mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
              {#if repositories.length > 0}
                {#each repositories as repo}
                  <button
                    type="button"
                    class="w-full text-left px-4 py-3 hover:bg-base-200 border-b border-base-300 last:border-b-0 focus:bg-base-200 focus:outline-none"
                    on:click={() => selectRepository(repo)}
                  >
                    <div class="font-medium text-base-content">{repo.name}</div>
                    <div class="text-xs text-base-content/50 font-mono mt-1 truncate text-primary">{repo.path}</div>
                  </button>
                {/each}
              {:else}
                <div class="px-4 py-6 text-center text-base-content/60">
                  <p class="font-medium">No matching repositories found</p>
                  <p class="text-sm text-base-content/50 mt-1">Try a different search term</p>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    </div>
    
    <!-- Sort Button -->
    <div class="flex items-end">
      <button class="btn btn-outline" on:click={toggleSort}>
        Sort 
        {#if sortDesc}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m-9 4h6" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m-6 4h6" />
          </svg>
        {/if}
      </button>
    </div>
  </div>
  
  <!-- Data Table -->
  <div class="flex-1 overflow-x-auto">
    <table class="table table-zebra w-full">
      <thead>
        <tr class="text-base-content bg-base-300">
          {#if isCompareMode}
            <th class="font-bold">Select</th>
          {/if}
          <th class="font-bold">ID</th>
          <th class="font-bold">Repository Name</th>
          <th class="font-bold">Date</th>
          <th class="font-bold">Status</th>
          <th class="font-bold">Vulnerabilities</th>
          {#if !isCompareMode}
            <th class="font-bold">Actions</th>
          {/if}
        </tr>
      </thead>
      <tbody>
        {#if isLoading}
          <tr>
            <td colspan="6" class="text-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
              <div class="mt-2">Loading scan history...</div>
            </td>
          </tr>
        {:else if scans.length === 0}
          <tr>
            <td colspan="{isCompareMode ? 6 : 6}" class="text-center text-base-content py-8">
              {#if selectedRepository}
                No scans found for repository "{selectedRepository.name}". 
              {:else if searchTerm}
                No scans found matching "{searchTerm}". Try a different search term.
              {:else}
                No scan history found. Try running a scan first.
              {/if}
            </td>
          </tr>
        {:else}
          {#each scans as scan, index (scan.scan.id)}
            <tr>
              {#if isCompareMode}
                <td class="text-center">
                  <input 
                    type="checkbox" 
                    class="checkbox checkbox-primary"
                    checked={selectedScansForComparison.some(s => s.scan.id === scan.scan.id)}
                    on:change={(e) => handleScanSelection(scan, e)}
                  />
                </td>
              {/if}
              <td class="font-medium text-base-content">{(currentPage - 1) * perPage + index + 1}</td>
              <td class="font-medium text-base-content">{scan.repository_name}</td>
              <td class="text-base-content">{formatDate(scan.scan.scan_time)}</td>
              <td>
                <span class="badge {getBadgeClass(scan.scan.status)} font-medium">{scan.scan.status}</span>
              </td>
              <td>
                {#if (scan.scan.vulnerabilities || 0) > 0}
                  <span class="text-error font-bold">{scan.vulnerabilities_count || scan.scan.vulnerabilities}</span>
                {:else}
                  <span class="text-success font-bold">{scan.vulnerabilities_count || scan.vulnerabilities || 0}</span>
                {/if}
              </td>
              {#if !isCompareMode}
              <td>
                <div class="flex gap-2">
                  <button 
                    class="btn btn-sm btn-outline btn-primary" 
                    on:click={() => viewScanDetail(scan)}
                    title="View scan details"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                    Details
                  </button>
                  <button 
                    class="btn btn-sm btn-outline btn-error" 
                    on:click={() => viewDeleteConfirmationScanData(scan)}
                    title="Delete scan"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </button>
                </div>
              </td>
              {/if}
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  
  <!-- Pagination -->
  {#if totalPages > 1}
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-base-content">
        Showing {startItem} to {endItem} of {totalItems} entries
        {#if selectedRepository}
          for repository "{selectedRepository.name}"
        {/if}
      </div>
      
      <div class="join">
        <button 
          class="join-item btn btn-sm"
          class:btn-disabled={currentPage === 1}
          on:click={prevPage}
        >
          «
        </button>
        
        {#each Array.from({length: Math.min(5, totalPages)}, (_, i) => {
          const start = Math.max(1, currentPage - 2);
          const end = Math.min(totalPages, start + 4);
          return start + i;
        }).filter(page => page <= totalPages) as page}
          <button 
            class="join-item btn btn-sm"
            class:btn-active={page === currentPage}
            on:click={() => goToPage(page)}
          >
            {page}
          </button>
        {/each}
        
        <button 
          class="join-item btn btn-sm"
          class:btn-disabled={currentPage === totalPages}
          on:click={nextPage}
        >
          »
        </button>
      </div>
    </div>
  {/if}
</div>
  
<!-- Modals -->
{#if showDetailModal && selectedScan}
  <HistoryDetailModal scan={selectedScan} on:close={closeDetailModal} />
{/if}

{#if showComparisonModal && comparisonData}
  <ComparisonModal 
    show={showComparisonModal}
    {comparisonData}
    on:close={() => showComparisonModal = false}
    on:viewDetail={(event) => viewComparisonVulnerabilityDetail(event.detail)}
  />
{/if}

{#if showConfirmModal && selectedScan}
  <ConfirmationModal 
    on:confirm={deleteScanData} 
    on:cancel={closeConfirmModal} 
    TitleConfirmationModal={titleConfirmModal} 
    MessageConfirmationModal={confirmMessage} 
  />
{/if}

{#if showInfoModal}
  <InformationModal 
    TitleModal={infoModalTitle}
    MessageModal={infoModalMessage}
    on:cancel={() => showInfoModal = false}
  />
{/if}
