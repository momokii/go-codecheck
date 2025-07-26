<script>
  import { UpdateUserPassword, CompleteFirstSetup } from '../../wailsjs/go/main/App';
  import { models } from '../../wailsjs/go/models';
  import { authStore, authActions } from '../stores';
  import InformationModal from '../components/general/InformationModal.svelte';

  let newPassword = '';
  let confirmPassword = '';
  let isLoading = false;
  let errorMessage = '';
  let successMessage = '';
  let showErrorModal = false;
  let showSuccessModal = false;

  $: user = $authStore.user;

  async function handlePasswordChange() {
    errorMessage = '';
    successMessage = '';

    // Validation
    if (!newPassword.trim()) {
      errorMessage = 'Please enter a new password.';
      showErrorModal = true;
      return;
    }

    if (newPassword.length < 6) {
      errorMessage = 'Password must be at least 6 characters long.';
      showErrorModal = true;
      return;
    }

    if (newPassword !== confirmPassword) {
      errorMessage = 'Passwords do not match. Please ensure both fields are identical.';
      showErrorModal = true;
      return;
    }

    isLoading = true;

    try {
      const userUpdate = new models.UserUpdate({
        id: user.id,
        password: newPassword,
        previous_password: 'admin', // this value is the default value for first login
      });

      await UpdateUserPassword(userUpdate, true); // true indicates this is completion of setup
      
      successMessage = 'Password has been updated successfully! You will now be redirected to the main application.';
      
      // Update user's setup completion status
      await CompleteFirstSetup(user.id);
      
      // Show success modal first
      showSuccessModal = true;
      
      // Clear the password change flag after showing success
      setTimeout(() => {
        authActions.clearPasswordChangeFlag();
      }, 2500);

    } catch (error) {
      errorMessage = error || 'An unexpected error occurred.';
      showErrorModal = true;
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-orange-50 to-red-100 flex items-center justify-center p-4">
  <div class="bg-white rounded-xl shadow-2xl w-full max-w-md p-8">
    <!-- Header -->
    <div class="text-center mb-8">
      <div class="w-16 h-16 bg-warning rounded-xl flex items-center justify-center mx-auto mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-3a1 1 0 011-1h2.586l6.414-6.414a6 6 0 015.743-7.743z" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-gray-900">Set New Password</h1>
      <p class="text-gray-600 mt-2">Please set a new password to complete your setup</p>
    </div>

    <!-- Password Change Form -->
    <div class="relative">
      {#if isLoading}
        <div class="absolute inset-0 bg-white bg-opacity-75 rounded-lg flex items-center justify-center z-10">
          <div class="text-center">
            <div class="loading loading-spinner loading-lg text-warning mb-2"></div>
            <p class="text-sm text-gray-600">Updating password...</p>
          </div>
        </div>
      {/if}
      
      <form on:submit|preventDefault={handlePasswordChange} class="space-y-6" class:opacity-75={isLoading}>
        <!-- New Password Field -->
        <div class="form-control">
          <label class="label" for="newPassword">
            <span class="label-text font-medium text-gray-700">New Password</span>
          </label>
          <input
            id="newPassword"
            type="password"
            placeholder="Enter your new password"
            class="input input-bordered w-full focus:input-warning text-black"
            bind:value={newPassword}
            disabled={isLoading}
            required
            aria-describedby="password-help"
          />
          <div class="label">
            <span id="password-help" class="label-text-alt text-gray-500">Minimum 6 characters</span>
          </div>
        </div>

        <!-- Confirm Password Field -->
        <div class="form-control">
          <label class="label" for="confirmPassword">
            <span class="label-text font-medium text-gray-700">Confirm New Password</span>
          </label>
          <input
            id="confirmPassword"
            type="password"
            placeholder="Confirm your new password"
            class="input input-bordered w-full focus:input-warning text-black"
            bind:value={confirmPassword}
            disabled={isLoading}
            required
          />
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          class="btn btn-warning w-full relative overflow-hidden"
          disabled={isLoading}
        >
          {#if isLoading}
            <span class="loading loading-spinner loading-sm"></span>
            <span class="ml-2">Updating Password...</span>
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-3a1 1 0 011-1h2.586l6.414-6.414a6 6 0 015.743-7.743z" />
            </svg>
            Set New Password
          {/if}
        </button>
      </form>
    </div>

    <!-- Info -->
    <div class="text-center mt-8 pt-6 border-t border-gray-200">
      <p class="text-sm text-gray-500">
        This is a one-time setup. After setting your password, you'll be taken to the main application.
      </p>
    </div>
  </div>

  <!-- Error Modal -->
  {#if showErrorModal && errorMessage}
    <InformationModal
      TitleModal="Password Update Error"
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
</div>

<style>
  /* Custom styles for better visual appeal */
  .input:focus {
    border-color: rgb(251, 191, 36);
    box-shadow: 0 0 0 3px rgba(251, 191, 36, 0.1);
  }
</style>
