<script>
  import { createEventDispatcher } from 'svelte';
  import { UpdateRepo, CheckIfFolderOrFIleExists } from '../../../wailsjs/go/main/App';
  
  const dispatch = createEventDispatcher();
  
  export let isOpen = false;
  export let repository = null;
  
  let formData = {
    id: 0,
    user_id: 1,
    name: '',
    description: '',
    path: ''
  };
  
  let errors = {};
  let isLoading = false;
  let isValidPath = false;
  
  // Reactive statement to validate path when it changes
  $: validatePath(formData.path);
  
  // Initialize form data when modal opens with repository data
  $: if (isOpen && repository && typeof repository === 'object' && repository['id']) {
    initializeFormData();
  }
  
  function initializeFormData() {
    try {
      const repo = repository;
      if (repo && repo['id']) {
        formData = {
          id: repo['id'] || 0,
          user_id: repo['user_id'] || 1,
          name: repo['name'] || '',
          description: repo['description'] || '',
          path: repo['path'] || ''
        };
        
        // Validate path when repository data loads
        if (repo['path']) {
          validatePath(repo['path']);
        }
      }
    } catch (e) {
      console.error('Error setting form data:', e);
    }
  }
  
  function closeModal() {
    isOpen = false;
    dispatch('close');
    resetForm();
  }
  
  function resetForm() {
    formData = {
      id: 0,
      user_id: 1,
      name: '',
      description: '',
      path: ''
    };
    errors = {};
    isValidPath = false;
  }
  
  async function validatePath(path) {
    if (!path || !path.trim()) {
      isValidPath = false;
      return;
    }
    
    try {
      isValidPath = await CheckIfFolderOrFIleExists(path.trim());
    } catch (error) {
      isValidPath = false;
    }
  }
  
  function validateForm() {
    errors = {};
    
    if (!formData.name.trim()) {
      errors.name = 'Repository name is required';
    }
    
    if (!formData.description.trim()) {
      errors.description = 'Description is required';
    }
    
    if (!formData.path.trim()) {
      errors.path = 'Path is required';
    } else if (!isValidPath) {
      errors.path = 'Path does not exist or is not accessible';
    }
    
    return Object.keys(errors).length === 0;
  }
  
  async function handleSubmit() {
    if (!validateForm()) {
      return;
    }
    
    isLoading = true;
    
    try {
      await UpdateRepo(formData.id, formData);
      dispatch('success', {
        title: 'Success',
        message: 'Repository updated successfully!'
      });
      closeModal();
    } catch (error) {
      dispatch('error', {
        title: 'Error',
        message: `Failed to update repository: ${error.message || error}`
      });
    } finally {
      isLoading = false;
    }
  }
  
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      closeModal();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
  <div class="modal modal-open" style="z-index: 1000;">
    <div class="modal-box max-w-md">
      <h3 class="font-bold text-lg mb-6 text-primary">Edit Repository</h3>
      
      <form on:submit|preventDefault={handleSubmit}>
        <!-- Repository Name -->
        <div class="form-control w-full mb-4">
          <label class="label" for="edit-repo-name">
            <span class="label-text font-medium">Repository Name</span>
          </label>
          <input
            id="edit-repo-name"
            type="text"
            bind:value={formData.name}
            placeholder="Enter repository name"
            class="input input-bordered w-full text-primary"
            class:input-error={errors.name}
            disabled={isLoading}
          />
          {#if errors.name}
            <div class="label">
              <span class="label-text-alt text-error">{errors.name}</span>
            </div>
          {/if}
        </div>
        
        <!-- Description -->
        <div class="form-control w-full mb-4">
          <label class="label" for="edit-repo-description">
            <span class="label-text font-medium">Description</span>
          </label>
          <textarea
            id="edit-repo-description"
            bind:value={formData.description}
            placeholder="Enter repository description"
            class="textarea textarea-bordered w-full text-primary"
            class:textarea-error={errors.description}
            rows="3"
            disabled={isLoading}
          ></textarea>
          {#if errors.description}
            <div class="label">
              <span class="label-text-alt text-error">{errors.description}</span>
            </div>
          {/if}
        </div>
        
        <!-- Path -->
        <div class="form-control w-full mb-6">
          <label class="label" for="edit-repo-path">
            <span class="label-text font-medium">Path</span>
          </label>
          <input
            id="edit-repo-path"
            type="text"
            bind:value={formData.path}
            placeholder="Enter repository path"
            class="input input-bordered w-full text-primary"
            class:input-error={errors.path}
            class:input-success={formData.path.trim() && isValidPath}
            disabled={isLoading}
          />
          {#if errors.path}
            <div class="label">
              <span class="label-text-alt text-error">{errors.path}</span>
            </div>
          {:else if formData.path.trim() && isValidPath}
            <div class="label">
              <span class="label-text-alt text-success">✓ Path is valid</span>
            </div>
          {:else if formData.path.trim() && !isValidPath}
            <div class="label">
              <span class="label-text-alt text-error">⚠ Path is Not Valid!</span>
            </div>
          {/if}
        </div>
        
        <!-- Action Buttons -->
        <div class="modal-action">
          <button
            type="button"
            class="btn btn-outline btn-error"
            on:click={closeModal}
            disabled={isLoading}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            disabled={isLoading || !formData.name.trim() || !formData.description.trim() || !formData.path.trim() || !isValidPath}
          >
            {#if isLoading}
              <span class="loading loading-spinner loading-sm"></span>
              Updating...
            {:else}
              Update Repository
            {/if}
          </button>
        </div>
      </form>
    </div>
    
    <!-- Modal backdrop -->
    <div class="modal-backdrop" on:click={closeModal} on:keydown={handleKeydown} role="button" tabindex="-1"></div>
  </div>
{/if}
