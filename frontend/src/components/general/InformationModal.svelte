<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();
  export let TitleModal = 'Information';
  export let MessageModal = 'Information Body';

  function cancel() {
    dispatch('cancel');
  }
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center" style="z-index: 2000;">
  <div class="bg-base-100 p-8 rounded-lg shadow-xl text-center max-w-lg">
    <h3 class="text-xl text-black font-bold mb-2">{TitleModal}</h3>
    <p class="mb-6 text-black">{MessageModal}</p>
    
    <div class="flex justify-center space-x-4">
      <button class="btn btn-primary" on:click={cancel}>Close</button>
    </div>
  </div>
</div>
