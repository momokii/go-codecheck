<script>
  import { createEventDispatcher } from 'svelte';
  
  export let show = false;
  export let comparisonData = {};
  
  const dispatch = createEventDispatcher();
  
  let activeTab = 'summary';
  let activeVulnTab = 'new';
  
  function close() {
    dispatch('close');
  }
  
  function handleKeydown(event, callback) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      callback();
    }
  }
  
  function viewVulnerabilityDetail(vulnerability) {
    dispatch('viewDetail', vulnerability);
  }
  
  function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      hour12: true
    }).format(date);
  }
  
  function getSeverityClass(severity) {
    if (!severity) return 'badge-neutral';
    
    switch(severity.toLowerCase()) {
      case 'critical':
        return 'badge-error text-white';
      case 'high':
        return 'badge-warning text-white';
      case 'medium':
        return 'badge-info text-white';
      case 'low':
        return 'badge-success text-white';
      default:
        return 'badge-neutral';
    }
  }
  
  function getStatusClass(status) {
    switch(status) {
      case 'new':
        return 'text-error font-semibold';
      case 'fixed':
        return 'text-success font-semibold';
      case 'unresolved':
        return 'text-warning font-semibold';
      default:
        return 'text-neutral';
    }
  }
  
  function getStatusIcon(status) {
    switch(status) {
      case 'new':
        return '⚠️';
      case 'fixed':
        return '✅';
      case 'unresolved':
        return '🔄';
      default:
        return '';
    }
  }
  
  // Safely render text data by ensuring proper string format and escaping
  function sanitizeText(text) {
    if (!text) return '';
    
    // Convert to string if it's not already
    let textStr = typeof text === 'string' ? text : String(text);
    
    // Remove any null bytes or invalid characters that could cause parsing issues
    textStr = textStr.replace(/\0/g, '').trim();
    
    // Ensure we have valid content
    if (!textStr) return 'No content available';
    
    return textStr;
  }
</script>

{#if show && comparisonData}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="modal modal-open" style="z-index: 40;" on:click|self={close}>
    <div class="modal-box w-11/12 max-w-7xl max-h-screen overflow-hidden" role="dialog" aria-labelledby="modal-title">
      <!-- Modal Header -->
      <div class="flex justify-between items-center mb-6 border-b pb-4">
        <div>
          <h3 id="modal-title" class="font-bold text-2xl text-primary">Scan Results Comparison</h3>
          <p class="text-sm text-gray-600 mt-1">
            Comparing vulnerabilities between two scan results
          </p>
        </div>
        <button class="btn btn-sm btn-circle btn-ghost" on:click={close}>✕</button>
      </div>

      <!-- Scan Information Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <!-- Older Scan -->
        <div class="card bg-base-200 border border-base-300">
          <div class="card-body p-4">
            <h4 class="card-title text-lg text-primary">📊 Baseline Scan</h4>
            <div class="space-y-2 text-sm">
              <div><strong>Repository:</strong> {comparisonData?.metadata?.olderScan?.name || 'Unknown'}</div>
              <div><strong>Date:</strong> {formatDate(comparisonData?.metadata?.olderScan?.date)}</div>
              <div><strong>Total Vulnerabilities:</strong> 
                <span class="badge badge-neutral ml-2">
                  {comparisonData?.metadata?.olderScan?.totalVulnerabilities || 0}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Newer Scan -->
        <div class="card bg-base-200 border border-base-300">
          <div class="card-body p-4">
            <h4 class="card-title text-lg text-primary">🔍 Latest Scan</h4>
            <div class="space-y-2 text-sm">
              <div><strong>Repository:</strong> {comparisonData?.metadata?.newerScan?.name || 'Unknown'}</div>
              <div><strong>Date:</strong> {formatDate(comparisonData?.metadata?.newerScan?.date)}</div>
              <div><strong>Total Vulnerabilities:</strong>
                <span class="badge badge-neutral ml-2">
                  {comparisonData?.metadata?.newerScan?.totalVulnerabilities || 0}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab Navigation -->
      <div class="tabs tabs-bordered mb-4">
        <button 
          class="tab tab-lg {activeTab === 'summary' ? 'tab-active' : ''}"
          on:click={() => activeTab = 'summary'}
        >
          📈 Summary
        </button>
        <button 
          class="tab tab-lg {activeTab === 'details' ? 'tab-active' : ''}"
          on:click={() => activeTab = 'details'}
        >
          🔍 Detailed View
        </button>
      </div>

      <!-- Tab Content -->
      <div class="h-96 overflow-auto">
        {#if activeTab === 'summary'}
          <!-- Summary Tab -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- New Vulnerabilities -->
            <div class="card bg-error/10 border border-error/20">
              <div class="card-body p-6 text-center">
                <div class="text-4xl font-bold text-error mb-2">
                  {comparisonData?.summary?.new || 0}
                </div>
                <h4 class="text-lg font-semibold text-error mb-2">⚠️ New Vulnerabilities</h4>
                <p class="text-sm text-gray-600">
                  Issues found in the latest scan that weren't present in the baseline
                </p>
                {#if (comparisonData?.summary?.new || 0) > 0}
                  <div class="mt-4">
                    <div class="badge badge-error text-white">Action Required</div>
                  </div>
                {/if}
              </div>
            </div>

            <!-- Fixed Vulnerabilities -->
            <div class="card bg-success/10 border border-success/20">
              <div class="card-body p-6 text-center">
                <div class="text-4xl font-bold text-success mb-2">
                  {comparisonData?.summary?.fixed || 0}
                </div>
                <h4 class="text-lg font-semibold text-success mb-2">✅ Fixed Vulnerabilities</h4>
                <p class="text-sm text-gray-600">
                  Issues that were resolved since the baseline scan
                </p>
                {#if (comparisonData?.summary?.fixed || 0) > 0}
                  <div class="mt-4">
                    <div class="badge badge-success text-white">Great Progress!</div>
                  </div>
                {/if}
              </div>
            </div>

            <!-- Unresolved Vulnerabilities -->
            <div class="card bg-warning/10 border border-warning/20">
              <div class="card-body p-6 text-center">
                <div class="text-4xl font-bold text-warning mb-2">
                  {comparisonData?.summary?.unresolved || 0}
                </div>
                <h4 class="text-lg font-semibold text-warning mb-2">🔄 Unresolved Vulnerabilities</h4>
                <p class="text-sm text-gray-600">
                  Issues that persist from the baseline scan
                </p>
                {#if (comparisonData?.summary?.unresolved || 0) > 0}
                  <div class="mt-4">
                    <div class="badge badge-warning text-white">Needs Attention</div>
                  </div>
                {/if}
              </div>
            </div>
          </div>

          <!-- Summary Insights -->
          <div class="mt-8 card bg-base-200 border border-base-300">
            <div class="card-body p-6">
              <h4 class="text-lg font-semibold text-primary mb-4">🎯 Key Insights</h4>
              <div class="space-y-3">
                {#if (comparisonData?.summary?.new || 0) > 0}
                  <div class="alert alert-error">
                    <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.99-.833-2.74 0L3.073 19c-.77.833.192 2.5 1.732 2.5z" />
                    </svg>
                    <span><strong>{comparisonData?.summary?.new || 0}</strong> new vulnerabilities detected - immediate review recommended</span>
                  </div>
                {/if}
                
                {#if (comparisonData?.summary?.fixed || 0) > 0}
                  <div class="alert alert-success">
                    <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span><strong>{comparisonData?.summary?.fixed || 0}</strong> vulnerabilities have been successfully resolved</span>
                  </div>
                {/if}
                
                {#if (comparisonData?.summary?.unresolved || 0) > 0}
                  <div class="alert alert-warning">
                    <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.99-.833-2.74 0L3.073 19c-.77.833.192 2.5 1.732 2.5z" />
                    </svg>
                    <span><strong>{comparisonData?.summary?.unresolved || 0}</strong> vulnerabilities remain unresolved since the baseline</span>
                  </div>
                {/if}
              </div>
            </div>
          </div>
          
        {:else if activeTab === 'details'}
          <!-- Details Tab -->
          <div class="space-y-4">
            <!-- Vulnerability Category Tabs -->
            <div class="tabs tabs-boxed bg-base-200">
              <button 
                class="tab {activeVulnTab === 'new' ? 'tab-active' : ''}"
                on:click={() => activeVulnTab = 'new'}
              >
                ⚠️ New ({comparisonData?.summary?.new || 0})
              </button>
              <button 
                class="tab {activeVulnTab === 'fixed' ? 'tab-active' : ''}"
                on:click={() => activeVulnTab = 'fixed'}
              >
                ✅ Fixed ({comparisonData?.summary?.fixed || 0})
              </button>
              <button 
                class="tab {activeVulnTab === 'unresolved' ? 'tab-active' : ''}"
                on:click={() => activeVulnTab = 'unresolved'}
              >
                🔄 Unresolved ({comparisonData?.summary?.unresolved || 0})
              </button>
            </div>

            <!-- Vulnerability List -->
            <div class="overflow-auto max-h-80">
              {#if comparisonData && comparisonData.details && comparisonData.details[activeVulnTab] && comparisonData.details[activeVulnTab].length > 0}
                <div class="space-y-3">
                  {#each comparisonData.details[activeVulnTab] as vulnerability, index}
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <div class="card bg-base-100 border border-base-300 hover:border-primary cursor-pointer transition-colors"
                         role="button"
                         tabindex="0"
                         on:click={() => viewVulnerabilityDetail(vulnerability)}
                         on:keydown={(e) => handleKeydown(e, () => viewVulnerabilityDetail(vulnerability))}>
                      <div class="card-body p-4">
                        <div class="flex justify-between items-start">
                          <div class="flex-1">
                            <div class="flex items-center gap-3 mb-2">
                              <span class="text-lg">{getStatusIcon(vulnerability.status)}</span>
                              <h5 class="font-semibold text-primary hover:text-primary-focus">
                                {sanitizeText(vulnerability.check_id || 'Unknown Check')}
                              </h5>
                              <span class="badge {getSeverityClass(vulnerability.extra?.severity)}">
                                {vulnerability.extra?.severity || 'Unknown'}
                              </span>
                            </div>
                            
                            <p class="text-sm text-gray-600 mb-2 line-clamp-2">
                              {sanitizeText(vulnerability.extra?.message || 'No description available')}
                            </p>
                            
                            <div class="flex items-center gap-4 text-xs text-gray-500">
                              <span><strong>File:</strong> {sanitizeText(vulnerability.path || 'Unknown')}</span>
                              {#if vulnerability.start}
                                <span><strong>Line:</strong> {vulnerability.start.line}</span>
                              {/if}
                              <span class="{getStatusClass(vulnerability.status)}">
                                <strong>{vulnerability.status.toUpperCase()}</strong>
                              </span>
                            </div>
                          </div>
                          
                          <div class="flex-shrink-0 ml-4">
                            <button class="btn btn-sm btn-outline btn-primary">
                              View Details
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              {:else}
                <div class="text-center py-12">
                  <div class="text-6xl mb-4">
                    {#if activeVulnTab === 'new'}🎉{:else if activeVulnTab === 'fixed'}✨{:else}👍{/if}
                  </div>
                  <h4 class="text-xl font-semibold text-gray-600 mb-2">
                    No {activeVulnTab} vulnerabilities
                  </h4>
                  <p class="text-gray-500">
                    {#if activeVulnTab === 'new'}
                      Great! No new security issues were introduced.
                    {:else if activeVulnTab === 'fixed'}
                      No vulnerabilities were fixed in this comparison period.
                    {:else}
                      Excellent! No persistent vulnerabilities found.
                    {/if}
                  </p>
                </div>
              {/if}
            </div>
          </div>
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="modal-action border-t pt-4 mt-4">
        <button class="btn btn-primary" on:click={close}>
          Close Comparison
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
