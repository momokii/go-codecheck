<script>
  import { createEventDispatcher } from 'svelte';
  import { isScanning, scanHistory, authStore } from '../../stores';
  
  import { 
    InitAndPrepareFolderScanSemgrep,
    RunSemgrepScan,
    GetSemgrepReportData,
    GetRepoDatas,
    CreateNewScan
  } from '../../../wailsjs/go/main/App';
  
  const dispatch = createEventDispatcher();
  
  $: user = $authStore.user;
  
  let searchTerm = '';
  let selectedRepository = null;
  let repositories = [];
  let showDropdown = false;
  let scanError = null;
  let scanStage = ''; // idle, preparing, scanning, processing
  let searchTimeout = null;
  let using_all_rules = true
  
  // Search for repositories as user types
  async function searchRepositories() {
    if (!searchTerm.trim()) {
      repositories = [];
      showDropdown = false;
      return;
    }
    
    if (!user?.id) {
      repositories = [];
      showDropdown = true;
      return;
    }
    
    try {
      const response = await GetRepoDatas(user.id, 1, 5, searchTerm, false);
      repositories = Array.isArray(response.data) ? response.data : [];
      // Show dropdown if we have results OR if we searched but found nothing
      showDropdown = true;
    } catch (error) {
      repositories = [];
      showDropdown = true; // Show dropdown to display error state
    }
  }
  
  // Handle search input with debouncing
  function handleSearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      searchRepositories();
    }, 300);
  }
  
  // Select a repository from the dropdown
  function selectRepository(repo) {
    selectedRepository = repo;
    searchTerm = repo.name;
    repositories = [];
    showDropdown = false;
  }
  
  // Clear selection
  function clearSelection() {
    selectedRepository = null;
    searchTerm = '';
    repositories = [];
    showDropdown = false;
  }
  
  async function startScan() {
    if (!selectedRepository) {
      scanError = "Please select a repository to scan";
      return;
    }
    
    const projectPath = selectedRepository.path;
    
    try {
      // Set scanning state
      $isScanning = true;
      scanError = null;
      
      // Stage 1: Prepare folder for scanning
      scanStage = 'preparing';
      await InitAndPrepareFolderScanSemgrep(projectPath);
      
      // Stage 2: Run the scan
      // Stage 2: Run the scan
      scanStage = 'scanning';
      await RunSemgrepScan(using_all_rules);
      
      // Stage 3: Process results
      scanStage = 'processing';
      const reportData = await GetSemgrepReportData();
      
      // Prepare scan data for database
      const scanData = {
        user_id: user?.id,
        repository_id: selectedRepository.id,
        repo_name: selectedRepository.name,
        repo_path: projectPath,
        vulnerabilities: reportData.results?.length || 0,
        scan_time: new Date().toISOString(),
        result: JSON.stringify(reportData), // Store the full report as JSON
        status: 'completed'
      };

      // Save scan to database
      await CreateNewScan(1, scanData); // userId first, then scan data

      // Add scan to local history (for immediate UI update)
      const newScan = {
        id: Date.now(), // Temporary ID for local state
        name: selectedRepository.name,
        path: projectPath,
        date: new Date().toISOString(),
        status: 'completed',
        vulnerabilities: reportData.results?.length || 0,
        report: reportData
      };
      
      $scanHistory = [newScan, ...$scanHistory];
      
      // Reset state and close modal
      $isScanning = false;
      dispatch('close');

      // Dispatch success event with details
      dispatch('scanComplete', {
        success: true,
        title: "Scan Success",
        message: `Scan of ${newScan.name} completed successfully. ${newScan.vulnerabilities} vulnerabilities found.`
      });
      
    } catch (error) {
      scanError = error || "Failed to complete scan 2";
      $isScanning = false;

      // Dispatch error event
      dispatch('scanComplete', {
        success: false,
        title: "Scan Failed",
        message: error || "Failed to complete scan 5"
      });
    }
  }
  
  async function cancelScan() {
    if ($isScanning) {
      // We would call a cancel API here if available
      $isScanning = false;
    }
    dispatch('close');
  }
  
  // Close dropdown when clicking outside
  function handleClickOutside(event) {
    if (!event.target.closest('.relative')) {
      showDropdown = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
  <div class="bg-base-100 rounded-lg shadow-xl w-full max-w-lg">    <!-- Modal Header -->
    <div class="p-4 border-b border-gray-200">
      <h3 class="text-xl font-bold text-primary">Scan Project</h3>
    </div>
    
    <!-- Modal Body -->
    <div class="p-6">
      {#if scanError}
        <div class="alert alert-error mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{scanError}</span>
        </div>
      {/if}
      
      {#if $isScanning}
        <!-- Scanning Progress -->
        <div class="flex flex-col items-center justify-center py-4">
          <div class="loading loading-spinner loading-lg"></div>
          <p class="mt-4 text-center text-bold text-primary">
            {#if scanStage === 'preparing'}
              Preparing project files for scanning...
            {:else if scanStage === 'scanning'}
              Running security scan on your code...
            {:else if scanStage === 'processing'}
              Processing scan results...
            {:else}
              Scanning in progress...
            {/if}
          </p>
        </div>      {:else}
        <!-- Repository Search -->
        <div class="form-control w-full">
          <label for="repository-search" class="label">
            <span class="label-text font-medium text-primary-focus">Search and Select Repository</span>
          </label>
          
          <!-- Search Input -->
          <div class="relative">
            <input 
              id="repository-search"
              type="text" 
              class="input input-bordered w-full text-base-content font-medium" 
              placeholder="Type repository name to search..." 
              bind:value={searchTerm}
              on:input={handleSearchInput}
              on:focus={() => searchTerm && searchRepositories()}
              class:input-success={selectedRepository} />
            
            <!-- Clear button -->
            {#if searchTerm}
              <button 
                type="button"
                class="absolute right-2 top-1/2 transform -translate-y-1/2 btn btn-ghost btn-sm btn-circle"
                on:click={clearSelection}
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            {/if}
            
            <!-- Search Results Dropdown -->
            {#if showDropdown}
              <div class="absolute top-full left-0 right-0 z-10 mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                {#if repositories.length > 0}
                  {#each repositories as repo}
                    <button
                      type="button"
                      class="w-full text-left px-4 py-3 hover:bg-base-200 border-b border-base-300 last:border-b-0 focus:bg-base-200 focus:outline-none"
                      on:click={() => selectRepository(repo)}
                    >
                      <div class="font-medium text-base-content">{repo.name}</div>
                      <div class="text-sm text-base-content/70 truncate text-primary">{repo.description}</div>
                      <div class="text-xs text-base-content/50 font-mono mt-1 truncate text-primary">{repo.path}</div>
                    </button>
                  {/each}
                {:else}
                  <!-- No results found state -->
                  <div class="px-4 py-6 text-center text-base-content/60">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 mx-auto mb-2 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <p class="font-medium">No matching repositories found</p>
                    <p class="text-sm text-base-content/50 mt-1">Try a different search term</p>
                  </div>
                {/if}
              </div>
            {/if}
          </div>
          
          <div class="label">
            <span class="label-text-alt font-medium text-neutral">Type to search for repositories in your database</span>
          </div>
        </div>
        
        <!-- Selected Repository Display -->
        {#if selectedRepository}
          <div class="mt-4">
            <label class="label" for="selected-repo-path">
              <span class="label-text font-medium text-primary-focus">Selected Repository Path</span>
            </label>
            <textarea 
              id="selected-repo-path"
              class="textarea textarea-bordered w-full font-mono text-sm text-black"
              rows="2"
              readonly
              value={selectedRepository.path}
            ></textarea>
          </div>
        {/if}

        <!-- CHOOSE USING ALL RULES SEMGREP OR DEFAULT SETUP -->
        <div class="mt-4">
          <label class="label" for="scan-setup-select">
            <span class="label-text font-medium text-primary-focus">Choose Scan Setup</span>
          </label>
          <select id="scan-setup-select" class="select select-bordered w-full text-black" bind:value={using_all_rules}>
            <option value={true}>Use All Rules (Recommended)</option>
            <option value={false}>Use Default Setup (Less Than All Rules)</option>
          </select>
        </div>

        <!-- Information Box -->
        <div class="mt-4 p-3 bg-base-200 rounded-lg shadow-sm">
          <p class="mb-2 font-large text-bold text-black">The selected repository will be:</p>
          <ol class="list-decimal pl-5 text-base-content">
            <li class="mb-1">Copied to a temporary location for scanning</li>
            <li class="mb-1">Analyzed for security vulnerabilities</li>
            <li class="mb-1">Results will be displayed in the application</li>
          </ol>
          <p class="mt-2 text-warning font-medium">Note: Your original files will not be modified</p>
        </div>
      {/if}
    </div>
    
    <!-- Modal Footer -->
    <div class="flex items-center justify-end p-4 border-t border-gray-200 gap-2">
      <button 
        class="btn btn-outline btn-error"
        on:click={cancelScan}
        disabled={$isScanning && scanStage !== 'scanning'}>
        {$isScanning ? 'Cancel' : 'Close'}
      </button>
      
      {#if !$isScanning}
        <button 
          class="btn btn-primary" 
          on:click={startScan}
          disabled={!selectedRepository}>
          Start Scan
        </button>
      {/if}
    </div>
  </div>
</div>
