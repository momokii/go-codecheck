<script>
  import { Login, GetAndValidateUserByToken } from '../../wailsjs/go/main/App';
  import { models } from '../../wailsjs/go/models';
  import { authActions } from '../stores';
  import InformationModal from '../components/general/InformationModal.svelte';
  let username = '';
  let password = '';
  let isLoading = false;
  let errorMessage = '';
  let showInfoModal = false;

  async function handleLogin() {
    if (!username.trim() || !password.trim()) {
      errorMessage = 'Please enter both username and password';
      return;
    }

    isLoading = true;
    errorMessage = '';

    try {
      const userLogin = new models.UserLogin({
        username: username.trim(),
        password: password
      });
      
      // First, login and get the token
      const token = await Login(userLogin);
      
      // Then, validate the token and get full user information
      const userData = await GetAndValidateUserByToken(token);
      
      // Create user object with complete info from backend
      const user = {
        id: userData.id,
        username: userData.username,
        is_completed_setup: userData.is_completed_setup
      };

      // Check if this is a first login (setup not completed)
      const needsPasswordChange = !userData.is_completed_setup;
      
      authActions.login(user, token, needsPasswordChange);
      
    } catch (error) {
      errorMessage = error
      showInfoModal = true;
    } finally {
      isLoading = false;
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
  <div class="bg-white rounded-xl shadow-2xl w-full max-w-md p-8">
    <!-- Logo/Header -->
    <div class="text-center mb-8">
      <div class="w-16 h-16 bg-primary rounded-xl flex items-center justify-center mx-auto mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-gray-900">Welcome Back</h1>
      <p class="text-gray-600 mt-2">Sign in to access CodeCheck Desktop</p>
    </div>

    <!-- Login Form -->
    <div class="relative">
      {#if isLoading}
        <div class="absolute inset-0 bg-white bg-opacity-75 rounded-lg flex items-center justify-center z-10">
          <div class="text-center">
            <div class="loading loading-spinner loading-lg text-primary mb-2"></div>
            <p class="text-sm text-gray-600">Authenticating...</p>
          </div>
        </div>
      {/if}
      
      <form on:submit|preventDefault={handleLogin} class="space-y-6" class:opacity-75={isLoading}>
      <!-- Username Field -->
      <div class="form-control">
        <label class="label" for="username">
          <span class="label-text font-medium text-gray-700">Username</span>
        </label>
        <input
          id="username"
          type="text"
          placeholder="Enter your username"            class="input input-bordered w-full focus:input-primary text-black"
            bind:value={username}
          on:keypress={handleKeyPress}
          disabled={isLoading}
          required
        />
      </div>

      <!-- Password Field -->
      <div class="form-control">
        <label class="label" for="password">
          <span class="label-text font-medium text-gray-700">Password</span>
        </label>
        <input
          id="password"
          type="password"
          placeholder="Enter your password"            class="input input-bordered w-full focus:input-primary text-black"
            bind:value={password}
          on:keypress={handleKeyPress}
          disabled={isLoading}
          required
        />        </div>

        <!-- Login Button -->
      <button
        type="submit"
        class="btn btn-primary w-full relative overflow-hidden"
        class:loading={isLoading}
        disabled={isLoading}
      >
        {#if isLoading}
          <span class="loading loading-spinner loading-sm"></span>
          <span class="ml-2">Signing in...</span>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
          </svg>
          Sign In
        {/if}        </button>
      </form>
    </div>

    <!-- Footer -->
    <div class="text-center mt-8 pt-6 border-t border-gray-200">
      <p class="text-sm text-gray-500">
        CodeCheck Desktop v1.0
      </p>
    </div>
  </div>

  <!-- Error Modal -->
  {#if showInfoModal && errorMessage}
    <InformationModal
      TitleModal="Login Error"
      MessageModal={errorMessage}
      on:cancel={() => {
        showInfoModal = false;
        errorMessage = '';
      }}
    />
  {/if}
</div>

<style>
  /* Custom styles for better visual appeal */
  .input:focus {
    border-color: rgb(79, 70, 229);
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
  }
</style>
