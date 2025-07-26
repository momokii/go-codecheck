<script>
  import { sidebarOpen, activeMenu, authStore, authActions } from '../../stores';
  import { Logout } from '../../../wailsjs/go/main/App';
  import UserManagementModal from './UserManagementModal.svelte';
  import logoImage from '../../assets/images/logo-universal.png'

  $: user = $authStore.user;

  let showUserManagementModal = false;
  let showUserDropdown = false;

  function setActiveMenu(menu) {
    $activeMenu = menu;
  }

  function toggleUserDropdown() {
    showUserDropdown = !showUserDropdown;
  }

  function openUserManagement() {
    showUserDropdown = false;
    showUserManagementModal = true;
  }

  async function handleLogout() {
    showUserDropdown = false;
    
    try {
      if (user?.id) {
        await Logout(user.id);
      }
    } catch (error) {
      console.error('Logout error:', error);
      // Continue with logout even if backend call fails
    } finally {
      authActions.logout();
    }
  }

  // Close dropdown when clicking outside
  function handleClickOutside(event) {
    if (!event.target.closest('.user-dropdown-container')) {
      showUserDropdown = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

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
  
  <!-- User Menu Section -->
  {#if $sidebarOpen}
    <div class="user-dropdown-container relative p-2 border-t border-gray-700">
      <button 
        on:click={toggleUserDropdown}
        class="flex items-center w-full p-3 rounded-md hover:bg-gray-700 transition-colors duration-200">
        
        <!-- User Avatar -->
        <div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
        
        <div class="ml-3 flex-1 text-left">
          <div class="text-sm font-medium text-white truncate">
            {user?.username || 'User'}
          </div>
          <div class="text-xs text-gray-400">Account Settings</div>
        </div>
        
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor"
             class:rotate-180={showUserDropdown}>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      <!-- User Dropdown Menu -->
      {#if showUserDropdown}
        <div class="absolute bottom-full left-2 right-2 mb-2 bg-base-100 rounded-lg shadow-lg border border-gray-200 py-2 z-50">
          <button 
            on:click={openUserManagement}
            class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Manage Account
          </button>
          
          <hr class="my-1 border-gray-200">
          
          <button 
            on:click={handleLogout}
            class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Logout
          </button>
        </div>
      {/if}
    </div>
  {:else}
    <!-- Compact user menu for collapsed sidebar -->
    <div class="user-dropdown-container relative p-2 border-t border-gray-700">
      <button 
        on:click={toggleUserDropdown}
        class="flex items-center justify-center w-full p-2 rounded-md hover:bg-gray-700 transition-colors duration-200">
        
        <div class="w-6 h-6 bg-primary rounded-full flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
      </button>

      <!-- Compact User Dropdown Menu -->
      {#if showUserDropdown}
        <div class="absolute bottom-full left-full ml-2 mb-2 bg-base-100 rounded-lg shadow-lg border border-gray-200 py-2 z-50 whitespace-nowrap">
          <div class="px-4 py-2 text-sm font-medium text-gray-900 border-b border-gray-200">
            {user?.username || 'User'}
          </div>
          
          <button 
            on:click={openUserManagement}
            class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Manage Account
          </button>
          
          <button 
            on:click={handleLogout}
            class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Logout
          </button>
        </div>
      {/if}
    </div>
  {/if}
  
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

<!-- User Management Modal -->
<UserManagementModal 
  show={showUserManagementModal} 
  on:close={() => showUserManagementModal = false} 
/>
