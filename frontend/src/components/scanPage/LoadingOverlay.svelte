<script>
  import { isScanning } from '../../stores';
  
  // Function to cancel scanning (would call an API if available)
  function cancelScan() {
    $isScanning = false;
  }
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50">
  <div class="bg-base-100 p-8 rounded-lg shadow-xl text-center max-w-lg">
    <div class="loading loading-spinner loading-lg mb-6"></div>
    
    <h3 class="text-xl font-bold mb-2 text-primary">Scanning in Progress</h3>
    <p class="mb-6 text-gray-600">Please wait while we analyze your code for security vulnerabilities...</p>
    
    <button class="btn btn-error" on:click={cancelScan}>
      Cancel Scan
    </button>
  </div>
</div>
