<script>
  import { onMount } from 'svelte';
  import { GetAndValidateUserByToken } from '../wailsjs/go/main/App';
  
  import Sidebar from './components/general/Sidebar.svelte';
  import ScanView from './views/ScanView.svelte';
  import HistoryView from './views/HistoryView.svelte';
  import RepositoryView from './views/RepositoryView.svelte';
  import LoginView from './views/LoginView.svelte';
  import PasswordChangeView from './views/PasswordChangeView.svelte';
  
  import { activeMenu, sidebarOpen, authStore, authActions } from './stores';
  
  $: isAuthenticated = $authStore.isAuthenticated;
  $: needsPasswordChange = $authStore.needsPasswordChange;
  $: user = $authStore.user;
  $: token = $authStore.token;
  
  let isValidatingSession = true;

  onMount(async () => {
    // Check if we have a stored token and validate it
    if (token) {
      try {
        const validatedUser = await GetAndValidateUserByToken(token);
        
        // Update user information with validated data
        authActions.updateUser({
          id: validatedUser.id,
          username: validatedUser.username,
          is_completed_setup: validatedUser.is_completed_setup
        });
        
        // Check if user needs to complete setup (first login)
        if (!validatedUser.is_completed_setup) {
          authActions.login(
            {
              id: validatedUser.id,
              username: validatedUser.username,
              is_completed_setup: validatedUser.is_completed_setup
            },
            token,
            true // needsPasswordChange = true
          );
        }
        
        isValidatingSession = false;
        
      } catch (error) {
        console.log('Session validation failed:', error);
        // Token is invalid, clear auth state
        authActions.logout();
        isValidatingSession = false;
      }
    } else {
      // No token stored
      isValidatingSession = false;
    }
  });

  // Show loading during session validation
  if (isValidatingSession) {
    // You could create a proper loading component here
  }
</script>

{#if isValidatingSession}
  <!-- Loading Screen -->
  <div class="min-h-screen bg-base-200 flex items-center justify-center">
    <div class="text-center">
      <div class="loading loading-spinner loading-lg text-primary mb-4"></div>
      <p class="text-lg">Validating session...</p>
    </div>
  </div>
{:else if !isAuthenticated}
  <!-- Login Screen -->
  <LoginView />
{:else if needsPasswordChange}
  <!-- Password Change Screen -->
  <PasswordChangeView />
{:else}
  <!-- Main Application -->
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
{/if}

<style>
  /* Global styles can be added here if needed */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: "Nunito", -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  }
</style>
