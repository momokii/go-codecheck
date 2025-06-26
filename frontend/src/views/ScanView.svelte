<script>
  import { isScanning } from '../stores';
  import ScanModal from '../components/scanPage/ScanModal.svelte';
  import LoadingOverlay from '../components/scanPage/LoadingOverlay.svelte';
  import InformationModal from '../components/general/InformationModal.svelte'
  
  import { CheckIfFolderOrFIleExists, CheckDockerIsAvailable, CheckDockerImagesIsAvailable } from '../../wailsjs/go/main/App';
  
  let showScanModal = false;
  let dockerStatus = null;
  let semgrepStatus = null;
  
  let informationModal = false
  let informationModalTitle = "Information"
  let informationModalMessage = "Message"

  async function checkDocker() {
    try {
      const version = await CheckDockerIsAvailable();
      dockerStatus = true;
      return version;
    } catch (error) {
      dockerStatus = false;
      return null;
    }
  }
  
  async function checkSemgrep() {
    try {
      const image = await CheckDockerImagesIsAvailable('semgrep/semgrep', 'latest');
      semgrepStatus = true;
      return image;
    } catch (error) {
      semgrepStatus = false;
      return null;
    }
  }
  
  function openScanModal() {
    showScanModal = true;
  }
  
  function closeScanModal() {
    showScanModal = false;
  }

  function handleScanComplete(event) {
    // Extract event details
    const { success, title, message } = event.detail;
    
    // Set the information modal properties
    informationModal = true;
    informationModalTitle = title;
    informationModalMessage = message;
    
  }
</script>

<div class="p-6 flex flex-col h-full">
  <h1 class="text-3xl font-bold mb-8 text-primary">Code Scanner</h1>
  
  <!-- Status Cards -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
    <!-- Docker Status Card -->
    <div class="card bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-base-content font-bold">Docker Status</h2>
        {#if dockerStatus === null}
          <span class="badge badge-warning">Not Checked</span>
          <p class="text-base-content">Check Docker is Availability.</p>
          <button class="btn btn-primary mt-4" on:click={checkDocker}>Check Docker</button>
        {:else if dockerStatus === true}
          <span class="badge badge-success">Available</span>
          <p class="text-base-content">Docker is installed and running.</p>
          <button class="btn btn-outline btn-success mt-4" on:click={checkDocker}>Re-check</button>
        {:else}
          <span class="badge badge-error">Not Available</span>
          <p class="text-base-content">Docker is not installed or not running.</p>
          <button class="btn btn-outline btn-error mt-4" on:click={checkDocker}>Re-check</button>
        {/if}
      </div>
    </div>
      <!-- Semgrep Status Card -->
    <div class="card bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-base-content font-bold">Semgrep Status</h2>
        {#if semgrepStatus === null}
          <span class="badge badge-warning">Not Checked</span>
          <p class="text-base-content">Check Semgrep is Availability</p>
          <button class="btn btn-primary mt-4" on:click={checkSemgrep}>Check Semgrep</button>
        {:else if semgrepStatus === true}
          <span class="badge badge-success">Available</span>
          <p class="text-base-content">Semgrep Docker image is available.</p>
          <button class="btn btn-outline btn-success mt-4" on:click={checkSemgrep}>Re-check</button>
        {:else}
          <span class="badge badge-error">Not Available</span>
          <p class="text-base-content">Semgrep Docker image is not available.</p>
          <button class="btn btn-outline btn-error mt-4" on:click={checkSemgrep}>Re-check</button>
        {/if}
      </div>
    </div>
  </div>
  
  <!-- Scan Button (centered) -->
  <div class="flex-1 flex flex-col items-center justify-center">    <button 
      class="btn btn-primary btn-lg w-64 h-16 text-xl shadow-lg hover:scale-105 transition-transform duration-200"
      on:click={openScanModal}
      disabled={$isScanning || !dockerStatus || !semgrepStatus}>
      {#if $isScanning}
        <span class="loading loading-spinner loading-md mr-2"></span>
      {/if}
      Scan Now
    </button>
    <p class="text-sm mt-4 font-medium text-primary">Click to enter a project directory path to scan</p>
    
    {#if !dockerStatus || !semgrepStatus}
      <div class="alert alert-warning max-w-md mt-6">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
        <div>
          <h3 class="font-bold">Scanning Not Available</h3>
          <div class="text-sm">
            {#if !dockerStatus}
              <p>Docker is not available. Please install and start Docker.</p>
            {/if}
            {#if !semgrepStatus}
              <p>Semgrep image is not available. The image will be pulled automatically during scan.</p>
            {/if}
          </div>
        </div>
      </div>
    {/if}
  </div>
  
  <!-- Scan Modal -->
  {#if showScanModal}
    <ScanModal 
    on:close={closeScanModal} 
    on:scanComplete={handleScanComplete}
    />
  {/if}
  
  <!-- Loading Overlay (when scanning) -->
  {#if $isScanning}
    <LoadingOverlay />
  {/if}

  <!-- information modal -->
  {#if informationModal}
    <InformationModal 
    TitleModal={informationModalTitle}
    MessageModal={informationModalMessage}
    on:cancel={() => informationModal = false} 
    />
  {/if}
</div>
