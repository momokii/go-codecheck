<script>
  import { createEventDispatcher } from 'svelte';
  import { isScanning, scanHistory } from '../../stores';
  
  import { 
    InitAndPrepareFolderScanSemgrep,
    RunSemgrepScan,
    GetSemgrepReportData,
    // CopyExternalProjectToScanDir,
    CheckIfFolderOrFIleExists
  } from '../../../wailsjs/go/main/App';
  
  const dispatch = createEventDispatcher();
  
  let projectPath = '';
  let scanError = null;
  let scanStage = ''; // idle, preparing, scanning, processing
  let isValidPath = false;
  
  async function validatePath() {
    if (!projectPath) {
      isValidPath = false;
      return;
    }
    
    try {
      isValidPath = await CheckIfFolderOrFIleExists(projectPath);
    } catch (error) {
      isValidPath = false;
    }
  }
  
  async function startScan() {
    if (!projectPath) {
      scanError = "Please enter a project directory path to scan";
      return;
    }
    
    try {
      // Set scanning state
      $isScanning = true;
      scanError = null;
      
      // Stage 1: Copy external project to scan directory
      scanStage = 'preparing';
      
      // Stage 2: Prepare folder for scanning
      await InitAndPrepareFolderScanSemgrep(projectPath);
      
      // Stage 3: Run the scan
      scanStage = 'scanning';
      const scanResult = await RunSemgrepScan();
      
      // Stage 4: Process results
      scanStage = 'processing';
      const reportData = await GetSemgrepReportData();
      
      // Add scan to history
      const newScan = {
        id: $scanHistory.length + 1,
        name: projectPath.split('/').pop().split('\\').pop(), // Extract folder name
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

      // Dispatch success event with details and this information for show the information modal
      dispatch('scanComplete', {
        success: true,
        title: "Scan Success",
        message: `Scan of ${newScan.name} completed successfully. ${newScan.vulnerabilities} vulnerabilities found.`
      });
      
      
    } catch (error) {
      scanError = error.message || "Failed to complete scan";
      $isScanning = false;

      // Dispatch error event, same for information modal if scan failed
      dispatch('scanComplete', {
        success: false,
        title: "Scan Failed",
        message: error.message || "Failed to complete scan"
      });
    }
  }
  
  function cancelScan() {
    if ($isScanning) {
      // We would call a cancel API here if available
      $isScanning = false;
    }
    dispatch('close');
  }
</script>

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
        <!-- Project Path Input -->
        <div class="form-control w-full">          <label for="project-path" class="label">
            <span class="label-text font-medium text-primary-focus">Enter project directory path to scan</span>
          </label>
          <div class="flex space-x-2">            <input 
              id="project-path"
              type="text" 
              class="input input-bordered flex-1 text-base-content font-medium" 
              placeholder="Enter absolute path to project directory" 
              bind:value={projectPath}
              on:input={validatePath}
              class:input-error={projectPath && !isValidPath}
              class:input-success={projectPath && isValidPath} />
          </div>
          <div class="label">
            <span class="label-text-alt font-medium text-neutral text-bold">Enter the full path to your project directory (e.g., C:\Projects\my-app)</span>
          </div>
          <div class="mt-2 p-3 bg-base-200 rounded-lg shadow-sm">
            <p class="mb-2 font-large text-bold text-black ">The code from the specified directory will be:</p>
            <ol class="list-decimal pl-5 text-base-content">
              <li class="mb-1">Copied to a temporary location for scanning</li>
              <li class="mb-1">Analyzed for security vulnerabilities</li>
              <li class="mb-1">Results will be displayed in the application</li>            </ol>
            <p class="mt-2 text-warning font-medium">Note: Your original files will not be modified</p>
          </div>
        </div>
      {/if}
    </div>
    
    <!-- Modal Footer -->
    <div class="flex items-center justify-end p-4 border-t border-gray-200 gap-2">
      <button 
        class="btn btn-outline"
        on:click={cancelScan}
        disabled={$isScanning && scanStage !== 'scanning'}>
        {$isScanning ? 'Cancel' : 'Close'}
      </button>
      
      {#if !$isScanning}
        <button 
          class="btn btn-primary" 
          on:click={startScan}
          disabled={!projectPath || !isValidPath}>
          Start Scan
        </button>
      {/if}
    </div>
  </div>
</div>
