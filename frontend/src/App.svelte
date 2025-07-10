<script>
  import Sidebar from './components/general/Sidebar.svelte';
  import ScanView from './views/ScanView.svelte';
  import HistoryView from './views/HistoryView.svelte';
  import RepositoryView from './views/RepositoryView.svelte';
  
  import { activeMenu, sidebarOpen } from './stores';
  
  // Initialize app
  import { onMount } from 'svelte';
  
  onMount(() => {
    // Initialize anything needed on startup
  });
</script>

<main class="flex h-screen bg-base-200 text-neutral-content overflow-hidden">
  <!-- Sidebar -->
  <div class="h-full">
    <Sidebar />
  </div>
  
  <!-- Main Content -->
  <div class="flex-1 overflow-y-auto">
    {#if $activeMenu === 'scan'}
      <ScanView />
    {:else if $activeMenu === 'history'}
      <HistoryView />
    {:else if $activeMenu === 'repository'}
      <RepositoryView />
    {/if}
  </div>
</main>

<style>
  /* Global styles can be added here if needed */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: "Nunito", -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  }
</style>
