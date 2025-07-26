<script>
  import { createEventDispatcher } from 'svelte';
  import { UpdateUserUsername, UpdateUserPassword } from '../../../wailsjs/go/main/App';
  import { models } from '../../../wailsjs/go/models';
  import { authStore, authActions } from '../../stores';
  import InformationModal from './InformationModal.svelte';

  export let show = false;

  const dispatch = createEventDispatcher();

  $: user = $authStore.user;

  let activeTab = 'username';
  let isLoading = false;
  let errorMessage = '';
  let successMessage = '';
  let showErrorModal = false;
  let showSuccessModal = false;

  // Username change
  let newUsername = '';

  // Password change
  let currentPassword = '';
  let newPassword = '';
  let confirmPassword = '';

  function close() {
    resetForm();
    dispatch('close');
  }

  function resetForm() {
    activeTab = 'username';
    isLoading = false;
    errorMessage = '';
    successMessage = '';
    showErrorModal = false;
    showSuccessModal = false;
    newUsername = '';
    currentPassword = '';
    newPassword = '';
    confirmPassword = '';
  }

  async function handleUsernameChange() {
    if (!newUsername.trim()) {
      errorMessage = 'Please enter a new username.';
      showErrorModal = true;
      return;
    }

    if (newUsername.trim() === user.username) {
      errorMessage = 'New username must be different from current username.';
      showErrorModal = true;
      return;
    }

    if (newUsername.trim().length < 5) {
      errorMessage = 'New username must be atleast 5 character with alphanumeric'
      showErrorModal = true
      return
    } 

    isLoading = true;
    errorMessage = '';
    successMessage = '';
    showErrorModal = false;
    showSuccessModal = false;

    try {
      const userUpdate = new models.UserUpdate({
        id: user.id,
        username: newUsername.trim()
      });

      await UpdateUserUsername(userUpdate);
      
      // Update user in store
      authActions.updateUser({ username: newUsername.trim() });
      
      successMessage = 'Username has been updated successfully!';
      showSuccessModal = true;
      newUsername = '';
      
    } catch (error) {
      errorMessage = error
      showErrorModal = true;
    } finally {
      isLoading = false;
    }
  }

  async function handlePasswordChange() {
    errorMessage = '';
    successMessage = '';
    showErrorModal = false;
    showSuccessModal = false;

    // Validation
    if (!currentPassword.trim()) {
      errorMessage = 'Please enter your current password.';
      showErrorModal = true;
      return;
    }

    if (!newPassword.trim()) {
      errorMessage = 'Please enter a new password.';
      showErrorModal = true;
      return;
    }

    if (newPassword.length < 6) {
      errorMessage = 'New password must be at least 6 characters long.';
      showErrorModal = true;
      return;
    }

    if (newPassword !== confirmPassword) {
      errorMessage = 'New passwords do not match. Please ensure both fields are identical.';
      showErrorModal = true;
      return;
    }

    if (newPassword === currentPassword) {
      errorMessage = 'New password must be different from current password.';
      showErrorModal = true;
      return;
    }

    isLoading = true;

    try {
      const userUpdate = new models.UserUpdate({
        id: user.id,
        password: newPassword,
        previous_password: currentPassword // send current password for validation
      });

      await UpdateUserPassword(userUpdate, false); // false indicates this is not initial setup
      
      successMessage = 'Password has been updated successfully!';
      showSuccessModal = true;
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
      
    } catch (error) {
      errorMessage = error || 'An unexpected error occurred.';
      showErrorModal = true;
    } finally {
      isLoading = false;
    }
  }

  // Reset messages when switching tabs
  $: if (activeTab) {
    errorMessage = '';
    successMessage = '';
    showErrorModal = false;
    showSuccessModal = false;
  }
</script>

{#if show}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="modal modal-open" on:click|self={close}>
    <div class="modal-box w-11/12 max-w-lg" role="dialog" aria-labelledby="modal-title">
      <!-- Modal Header -->
      <div class="flex justify-between items-center mb-6">
        <h3 id="modal-title" class="font-bold text-xl text-primary">Account Settings</h3>
        <button class="btn btn-sm btn-circle btn-ghost" on:click={close}>✕</button>
      </div>

      <!-- Tab Navigation -->
      <div class="tabs tabs-bordered mb-6">
        <button 
          class="tab tab-lg {activeTab === 'username' ? 'tab-active' : ''}"
          on:click={() => activeTab = 'username'}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          Username
        </button>
        <button 
          class="tab tab-lg {activeTab === 'password' ? 'tab-active' : ''}"
          on:click={() => activeTab = 'password'}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-3a1 1 0 011-1h2.586l6.414-6.414a6 6 0 015.743-7.743z" />
          </svg>
          Password
        </button>
      </div>

      <!-- Tab Content -->
      <div class="min-h-64">
        {#if activeTab === 'username'}
          <!-- Username Change -->
          <div class="relative">
            {#if isLoading}
              <div class="absolute inset-0 bg-white bg-opacity-75 rounded-lg flex items-center justify-center z-10">
                <div class="text-center">
                  <div class="loading loading-spinner loading-lg text-primary mb-2"></div>
                  <p class="text-sm text-gray-600">Updating username...</p>
                </div>
              </div>
            {/if}
            
            <div class="space-y-4" class:opacity-75={isLoading}>
              <div class="form-control">
                <label class="label" for="currentUsername">
                  <span class="label-text font-medium">Current Username</span>
                </label>
                <input
                  id="currentUsername"
                  type="text"
                  value={user?.username || ''}
                  class="input input-bordered text-black"
                  disabled
                />
              </div>

              <div class="form-control">
                <label class="label" for="newUsername">
                  <span class="label-text font-medium">New Username</span>
                </label>
                <input
                  id="newUsername"
                  type="text"
                  placeholder="Enter new username"
                  class="input input-bordered focus:input-primary text-black"
                  bind:value={newUsername}
                  disabled={isLoading}
                />
              </div>

              <button
                class="btn btn-primary w-full relative overflow-hidden"
                disabled={isLoading || !newUsername.trim()}
                on:click={handleUsernameChange}
              >
                {#if isLoading}
                  <span class="loading loading-spinner loading-sm"></span>
                  <span class="ml-2">Updating...</span>
                {:else}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  Update Username
                {/if}
              </button>
            </div>
          </div>

        {:else if activeTab === 'password'}
          <!-- Password Change -->
          <div class="relative">
            {#if isLoading}
              <div class="absolute inset-0 bg-white bg-opacity-75 rounded-lg flex items-center justify-center z-10">
                <div class="text-center">
                  <div class="loading loading-spinner loading-lg text-primary mb-2"></div>
                  <p class="text-sm text-gray-600">Updating password...</p>
                </div>
              </div>
            {/if}
            
            <div class="space-y-4" class:opacity-75={isLoading}>
              <div class="form-control">
                <label class="label" for="currentPassword">
                  <span class="label-text font-medium">Current Password</span>
                </label>
                <input
                  id="currentPassword"
                  type="password"
                  placeholder="Enter current password"
                  class="input input-bordered focus:input-primary text-black"
                  bind:value={currentPassword}
                  disabled={isLoading}
                />
              </div>

              <div class="form-control">
                <label class="label" for="newPasswordModal">
                  <span class="label-text font-medium">New Password</span>
                </label>
                <input
                  id="newPasswordModal"
                  type="password"
                  placeholder="Enter new password"
                  class="input input-bordered focus:input-primary text-black"
                  bind:value={newPassword}
                  disabled={isLoading}
                  aria-describedby="password-help-modal"
                />
                <div class="label">
                  <span id="password-help-modal" class="label-text-alt text-gray-500">Minimum 6 characters</span>
                </div>
              </div>

              <div class="form-control">
                <label class="label" for="confirmPasswordModal">
                  <span class="label-text font-medium">Confirm New Password</span>
                </label>
                <input
                  id="confirmPasswordModal"
                  type="password"
                  placeholder="Confirm new password"
                  class="input input-bordered focus:input-primary text-black"
                  bind:value={confirmPassword}
                  disabled={isLoading}
                />
              </div>

              <button
                class="btn btn-primary w-full relative overflow-hidden"
                disabled={isLoading || !currentPassword.trim() || !newPassword.trim() || !confirmPassword.trim()}
                on:click={handlePasswordChange}
              >
                {#if isLoading}
                  <span class="loading loading-spinner loading-sm"></span>
                  <span class="ml-2">Updating...</span>
                {:else}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-3a1 1 0 011-1h2.586l6.414-6.414a6 6 0 015.743-7.743z" />
                  </svg>
                  Update Password
                {/if}
              </button>
            </div>
          </div>
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="modal-action">
        <button class="btn btn-outline" on:click={close}>Close</button>
      </div>
    </div>
  </div>

  <!-- Error Modal -->
  {#if showErrorModal && errorMessage}
    <InformationModal
      TitleModal="Error"
      MessageModal={errorMessage}
      on:cancel={() => {
        showErrorModal = false;
        errorMessage = '';
      }}
    />
  {/if}

  <!-- Success Modal -->
  {#if showSuccessModal && successMessage}
    <InformationModal
      TitleModal="Success"
      MessageModal={successMessage}
      on:cancel={() => {
        showSuccessModal = false;
        successMessage = '';
      }}
    />
  {/if}
{/if}
