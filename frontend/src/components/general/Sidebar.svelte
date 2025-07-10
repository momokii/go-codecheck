<script>
  import { sidebarOpen, activeMenu } from '../../stores';
  import logoImage from '../../assets/images/logo-universal.png'

  function setActiveMenu(menu) {
    $activeMenu = menu;
  }
</script>

<div class="h-full flex flex-col bg-neutral shadow-lg transition-all duration-300" 
     class:w-64={$sidebarOpen} 
     class:w-16={!$sidebarOpen}>
  
  <!-- Logo and App Title -->
  <div class="flex items-center p-4 border-b border-gray-700">
    <div class="w-8 h-8 mr-3 rounded-md bg-primary flex items-center justify-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-primary-content" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    </div>
    {#if $sidebarOpen}
      <h1 class="text-lg font-semibold text-white">CodeCheck Desktop</h1>
    {/if}
  </div>
  
  <!-- Menu Items -->
  <nav class="flex flex-col gap-2 p-2 flex-1">
    <!-- Scan Menu Button -->
    <button 
      on:click={() => setActiveMenu('scan')}
      class="flex items-center p-3 rounded-md transition-colors duration-200"
      class:bg-primary={$activeMenu === 'scan'}
      class:bg-opacity-20={$activeMenu === 'scan'}
      class:hover:bg-gray-700={$activeMenu !== 'scan'}>
      
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      
      {#if $sidebarOpen}
        <span class="ml-3">Scan</span>
      {/if}
    </button>
    
    <!-- History Menu Button -->
    <button 
      on:click={() => setActiveMenu('history')}
      class="flex items-center p-3 rounded-md transition-colors duration-200"
      class:bg-primary={$activeMenu === 'history'}
      class:bg-opacity-20={$activeMenu === 'history'}
      class:hover:bg-gray-700={$activeMenu !== 'history'}>
      
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      
      {#if $sidebarOpen}
        <span class="ml-3">History</span>
      {/if}
    </button>

    <!-- Repo Menu Button -->
    <button 
      on:click={() => setActiveMenu('repository')}
      class="flex items-center p-3 rounded-md transition-colors duration-200"
      class:bg-primary={$activeMenu === 'repository'}
      class:bg-opacity-20={$activeMenu === 'repository'}
      class:hover:bg-gray-700={$activeMenu !== 'repository'}>
      
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      
      {#if $sidebarOpen}
        <span class="ml-3">Repository</span>
      {/if}
    </button>



  </nav>
  
  <!-- Toggle Sidebar Button -->
  <div class="p-4 border-t border-gray-700">
    <button 
      on:click={() => $sidebarOpen = !$sidebarOpen}
      class="flex items-center justify-center w-full p-2 rounded-md hover:bg-gray-700 transition-colors duration-200">
      
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"
           class:rotate-180={!$sidebarOpen}>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
      </svg>
      
      {#if $sidebarOpen}
        <span class="ml-2">Close</span>
      {/if}
    </button>
  </div>
</div>
