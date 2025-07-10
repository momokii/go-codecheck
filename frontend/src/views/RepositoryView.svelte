<script>
  import { onMount } from 'svelte';
  import InformationModal from '../components/general/InformationModal.svelte';
  import ConfirmationModal from '../components/general/ConfirmationModal.svelte';
  import AddRepositoryModal from '../components/repositoryPage/AddRepositoryModal.svelte';
  import EditRepositoryModal from '../components/repositoryPage/EditRepositoryModal.svelte';
  
  import { GetRepoDatas, DeleteRepo } from '../../wailsjs/go/main/App';
  
  // State variables
  let repositories = [];
  let isLoading = false;
  let currentPage = 1;
  let perPage = 10;
  let totalItems = 0;
  let totalPages = 0;
  let searchTerm = '';
  let sortDesc = false;
  
  // Modal states
  let showAddModal = false;
  let showEditModal = false;
  let showConfirmModal = false;
  let showInfoModal = false;
  
  // Modal data
  let selectedRepository = null;
  let infoModalTitle = '';
  let infoModalMessage = '';
  let confirmModalTitle = '';
  let confirmModalMessage = '';
  
  // Computed
  $: totalPages = Math.ceil(totalItems / perPage);
  $: startItem = (currentPage - 1) * perPage + 1;
  $: endItem = Math.min(currentPage * perPage, totalItems);
  
  onMount(() => {
    loadRepositories();
  });
  
  async function loadRepositories() {
    isLoading = true;
    try {
      const userId = 1; // You might want to get this from a store
      const response = await GetRepoDatas(userId, currentPage, perPage, searchTerm, sortDesc);
      
      // Handle different possible response formats
      if (Array.isArray(response)) {
        repositories = response;
        // Since we don't have total count from backend, estimate it
        // If we get less than perPage items, we're on the last page
        if (response.length < perPage && currentPage === 1) {
          totalItems = response.length;
        } else if (response.length < perPage) {
          totalItems = (currentPage - 1) * perPage + response.length;
        } else {
          // Estimate that there might be more pages
          totalItems = currentPage * perPage + 1;
        }
      } else if (response && typeof response === 'object') {
        // Handle if response is an object with data and total properties
        repositories = [];
        totalItems = 0;
        try {
          // Try to access properties safely
          if (response.data) repositories = response.data;
          if (response.total) totalItems = response.total;
        } catch (e) {
          // Fallback to empty array
          repositories = [];
          totalItems = 0;
        }
      } else {
        repositories = [];
        totalItems = 0;
      }
    } catch (error) {
      showInfo('Error', `Failed to load repositories: ${error.message || error}`);
      repositories = [];
      totalItems = 0;
    } finally {
      isLoading = false;
    }
  }
  
  // Pagination functions
  function goToPage(page) {
    if (page >= 1 && page <= totalPages && page !== currentPage) {
      currentPage = page;
      loadRepositories();
    }
  }
  
  function nextPage() {
    if (currentPage < totalPages) {
      currentPage++;
      loadRepositories();
    }
  }
  
  function prevPage() {
    if (currentPage > 1) {
      currentPage--;
      loadRepositories();
    }
  }
  
  // Search function
  function handleSearch() {
    currentPage = 1;
    loadRepositories();
  }
  
  // Sort function
  function toggleSort() {
    sortDesc = !sortDesc;
    loadRepositories();
  }
  
  // Modal functions
  function openAddModal() {
    showAddModal = true;
  }
  
  function openEditModal(repository) {
    selectedRepository = repository;
    showEditModal = true;
  }
  
  function openDeleteConfirmation(repository) {
    selectedRepository = repository;
    confirmModalTitle = `Delete Repository "${repository.name}"`;
    confirmModalMessage = `Are you sure you want to delete the repository "${repository.name}"? This action cannot be undone.`;
    showConfirmModal = true;
  }
  
  function showInfo(title, message) {
    infoModalTitle = title;
    infoModalMessage = message;
    showInfoModal = true;
  }
  
  // Event handlers
  function handleAddSuccess(event) {
    showInfo(event.detail.title, event.detail.message);
    loadRepositories();
  }
  
  function handleAddError(event) {
    showInfo(event.detail.title, event.detail.message);
  }
  
  function handleEditSuccess(event) {
    showInfo(event.detail.title, event.detail.message);
    loadRepositories();
  }
  
  function handleEditError(event) {
    showInfo(event.detail.title, event.detail.message);
  }
  
  async function handleDeleteConfirm() {
    if (!selectedRepository) return;
    
    try {
      await DeleteRepo(selectedRepository.id, 1); // userId = 1
      showInfo('Success', 'Repository deleted successfully!');
      loadRepositories();
    } catch (error) {
      showInfo('Error', `Failed to delete repository: ${error.message || error}`);
    }
    
    showConfirmModal = false;
    selectedRepository = null;
  }
  
  function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric'
    }).format(date);
  }
</script>

<div class="p-6 flex flex-col h-full">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-3xl font-bold text-primary">Repository Management</h1>
    <button class="btn btn-primary" on:click={openAddModal}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Repository
    </button>
  </div>
  
  <!-- Search and Sort Controls -->
  <div class="flex gap-4 mb-6">
    <div class="flex-1">
      <div class="join">
        <input 
          type="text" 
          placeholder="Search repositories..." 
          class="join-item input input-bordered w-full max-w-md text-black"
          bind:value={searchTerm}
          on:keydown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button class="join-item btn btn-primary" on:click={handleSearch}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>
      </div>
    </div>
    
    <button class="btn btn-outline" on:click={toggleSort}>
      Sort 
      {#if sortDesc}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m-9 4h6" />
        </svg>
      {:else}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m-6 4h6" />
        </svg>
      {/if}
    </button>
  </div>
  
  <!-- Data Table -->
  <div class="flex-1 overflow-x-auto">
    <table class="table table-zebra w-full">
      <thead>
        <tr class="text-base-content bg-base-300">
          <th class="font-bold">ID</th>
          <th class="font-bold">Name</th>
          <th class="font-bold">Description</th>
          <th class="font-bold">Path</th>
          <th class="font-bold">Created</th>
          <th class="font-bold">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if isLoading}
          <tr>
            <td colspan="6" class="text-center py-8">
              <span class="loading loading-spinner loading-lg"></span>
              <div class="mt-2">Loading repositories...</div>
            </td>
          </tr>
        {:else if repositories.length === 0}
          <tr>
            <td colspan="6" class="text-center text-base-content py-8">
              {#if searchTerm}
                No repositories found matching "{searchTerm}". Try a different search term.
              {:else}
                No repositories found. Click "Add Repository" to create your first repository.
              {/if}
            </td>
          </tr>
        {:else}
          {#each repositories as repo, index (repo.id)}
            <tr>
              <td class="font-medium text-base-content">{(currentPage - 1) * perPage + index + 1}</td>
              <td class="font-medium text-base-content">{repo.name}</td>
              <td class="text-base-content">
                <div class="max-w-xs truncate" title={repo.description}>
                  {repo.description}
                </div>
              </td>
              <td class="text-base-content">
                <div class="max-w-xs truncate font-mono text-sm" title={repo.path}>
                  {repo.path}
                </div>
              </td>
              <td class="text-base-content">{formatDate(repo.create_at)}</td>
              <td>
                <div class="flex gap-2">
                  <button 
                    class="btn btn-sm btn-outline btn-primary"
                    on:click={() => openEditModal(repo)}
                    title="Edit repository"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button 
                    class="btn btn-sm btn-outline btn-error"
                    on:click={() => openDeleteConfirmation(repo)}
                    title="Delete repository"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  
  <!-- Pagination -->
  {#if totalPages > 1}
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-base-content">
        Showing {startItem} to {endItem} of {totalItems} entries
      </div>
      
      <div class="join">
        <button 
          class="join-item btn btn-sm"
          class:btn-disabled={currentPage === 1}
          on:click={prevPage}
        >
          «
        </button>
        
        {#each Array.from({length: Math.min(5, totalPages)}, (_, i) => {
          const start = Math.max(1, currentPage - 2);
          const end = Math.min(totalPages, start + 4);
          return start + i;
        }).filter(page => page <= totalPages) as page}
          <button 
            class="join-item btn btn-sm"
            class:btn-active={page === currentPage}
            on:click={() => goToPage(page)}
          >
            {page}
          </button>
        {/each}
        
        <button 
          class="join-item btn btn-sm"
          class:btn-disabled={currentPage === totalPages}
          on:click={nextPage}
        >
          »
        </button>
      </div>
    </div>
  {/if}
</div>

<!-- Modals -->
<AddRepositoryModal 
  bind:isOpen={showAddModal}
  on:success={handleAddSuccess}
  on:error={handleAddError}
  on:close={() => showAddModal = false}
/>

<EditRepositoryModal 
  bind:isOpen={showEditModal}
  bind:repository={selectedRepository}
  on:success={handleEditSuccess}
  on:error={handleEditError}
  on:close={() => showEditModal = false}
/>

{#if showConfirmModal}
  <ConfirmationModal 
    TitleConfirmationModal={confirmModalTitle}
    MessageConfirmationModal={confirmModalMessage}
    on:confirm={handleDeleteConfirm}
    on:cancel={() => showConfirmModal = false}
  />
{/if}

{#if showInfoModal}
  <InformationModal 
    TitleModal={infoModalTitle}
    MessageModal={infoModalMessage}
    on:cancel={() => showInfoModal = false}
  />
{/if}