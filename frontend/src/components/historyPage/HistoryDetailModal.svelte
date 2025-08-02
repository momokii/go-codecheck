<script>
  import { createEventDispatcher } from 'svelte';
  
  export let scan = {}; // The scan data passed from parent

  const dispatch = createEventDispatcher();
  
  // State for severity filter
  let selectedSeverity = 'ALL';
  
  // Reactive variable for filtered vulnerabilities
  $: filteredVulnerabilities = (() => {
    if (!scan.report?.results) return [];
    
    if (selectedSeverity === 'ALL') {
      return scan.report.results;
    } else {
      return scan.report.results.filter(
        result => result.extra?.severity?.toUpperCase() === selectedSeverity
      );
    }
  })();
  
  function closeModal() {
    dispatch('close');
  }
  
  function formatDate(dateString) {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric'
    }).format(date);
  }
  
  // Get severity class for badges
  function getSeverityClass(severity) {
    switch(severity.toLowerCase()) {
      case 'critical':
        return 'badge-error text-black';
      case 'error':
        return 'badge-error text-black';
      case 'warning':
        return 'badge-warning';
      case 'info':
        return 'badge-info';
      default:
        return 'badge-ghost';
    }
  }

  function getTotalVulnerabilities(scan) {
    if (!scan || !scan.report) return 0
    
    let vulnerabilities = 0

    const total_result = scan.report.results ? scan.report.results.length : 0

    if (total_result > 0 ) {
      for (const issue of scan.report.results) {
        if (issue.extra && issue.extra.severity && issue.extra.severity.toLowerCase() !== 'info') {
          vulnerabilities += 1
        }
      }
    }
    
    return vulnerabilities
  }

  
  // Safely render code snippets by escaping HTML and ensuring proper string format
  function sanitizeCodeSnippet(code) {
    if (!code) return '';
    
    // Convert to string if it's not already
    let codeStr = typeof code === 'string' ? code : String(code);
    
    // Remove any null bytes or invalid characters that could cause parsing issues
    codeStr = codeStr.replace(/\0/g, '').trim();
    
    // Ensure we have valid content
    if (!codeStr) return 'No code content available';
    
    return codeStr;
  }
  
  // Safely render text data by ensuring proper string format
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

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60] p-4 overflow-y-auto">
  <div class="bg-base-100 rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col">    <!-- Modal Header -->
    <div class="p-4 border-b border-gray-200 flex justify-between items-center">
      <h3 class="text-xl font-bold text-primary">
        Scan Report: {scan.name}
      </h3>
      <button class="btn btn-sm btn-circle" on:click={closeModal}>✕</button>
    </div>
    
    <!-- Modal Body -->
    <div class="p-6 overflow-y-auto flex-1">      <!-- Scan Summary -->
      <div class="card bg-base-200 mb-6 shadow-lg">
        <div class="card-body">
            <h4 class="card-title text-lg text-primary font-bold text-center w-full justify-center">Scan Summary</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p class="text-sm font-semibold text-primary-focus">Project Repositories</p>
              <p class="font-medium text-base-content">{scan.name}</p>
            </div>
            <div>
              <p class="text-sm font-semibold text-primary-focus">Project Path</p>
              <p class="font-medium text-base-content">{scan.path || 'N/A'}</p>
            </div>
            <div>
              <p class="text-sm font-semibold text-primary-focus">Scan Date</p>
              <p class="font-medium text-base-content">{formatDate(scan.date)}</p>
            </div>
            <div>
              <p class="text-sm font-semibold text-primary-focus">Status</p>
              <p>
                <span class="badge badge-{scan.status === 'completed' ? 'success' : 'warning'} text-white">
                  {scan.status}
                </span>
              </p>
            </div>            <div>
              <p class="text-sm font-semibold text-primary-focus">Total Issues</p>
              <p class="font-bold text-lg {scan.vulnerabilities  > 0 ? 'text-error' : 'text-success'}">
                {scan.vulnerabilities }
              </p>
            </div>
            <div>
              <p class="text-sm font-semibold text-primary-focus">Vulnerabilities</p>
              <p class="font-bold text-lg {getTotalVulnerabilities(scan) > 0 ? 'text-error' : 'text-success'}">
                {getTotalVulnerabilities(scan) }
              </p>
            </div>
            <div>
              <p class="text-sm font-semibold text-primary-focus">Informational Data</p>
              <p class="font-bold text-lg { getTotalVulnerabilities(scan) !== scan.vulnerabilities ? 'text-info' : 'text-success'}">
                {scan.vulnerabilities - getTotalVulnerabilities(scan) || 0}
              </p>
            </div>
            {#if scan.report && scan.report.version}
            <div>
              <p class="text-sm font-semibold text-primary-focus">Semgrep Version</p>
              <p class="font-medium text-base-content">{scan.report.version}</p>
            </div>
            {/if}
          </div>
        </div>
      </div>
      
      <!-- Severity Filter Dropdown -->
      {#if scan.report && scan.report.results && scan.report.results.length > 0}
        <div class="mb-6">
          <div class="flex items-center gap-3">
            <label for="severity-filter" class="text-sm font-semibold text-primary-focus">Filter by Severity:</label>
            <select 
              id="severity-filter"
              bind:value={selectedSeverity}
              class="select select-bordered select-sm w-auto min-w-[150px] bg-base-100 text-black"
            >
              <option value="ALL">All Severities</option>
              <option value="CRITICAL">Critical</option>
              <option value="ERROR">Error</option>
              <option value="WARNING">Warning</option>
              <option value="INFO">Info</option>
            </select>
            <span class="text-xs text-base-content opacity-70">
              Showing {filteredVulnerabilities.length} of {scan.report.results.length} vulnerabilities
            </span>
          </div>
        </div>
      {/if}

        <!-- Vulnerabilities List -->
      {#if scan.report && scan.report.results && scan.report.results.length > 0}
        <h4 class="text-lg font-bold mb-4 text-primary-focus">Vulnerabilities ({filteredVulnerabilities.length})</h4>
        
        <div class="space-y-6">
          {#each filteredVulnerabilities as issue}
            <div class="card border border-base-300 shadow-md bg-base-100">
              <div class="card-body p-4">
                <div class="flex justify-between items-start flex-wrap gap-2">
                  <h5 class="card-title text-base text-primary-focus">
                    {sanitizeText(issue.check_id || 'Unknown Issue')}
                  </h5>
                  <span class="badge {getSeverityClass(issue.extra?.severity || 'medium')} text-white font-medium">
                    {issue.extra?.severity || 'medium'}
                  </span>
                </div>
                
                <!-- Issue message -->
                <div class="my-2 text-base-content font-medium text-justify">
                  {sanitizeText(issue.extra?.message || 'No description available')}
                </div>
                
                <!-- Code snippet -->
                <!-- will not use this for now, because the code snipped required login to see it  -->
                <!-- {#if issue.extra?.lines}
                  <div class="bg-neutral text-neutral-content p-3 rounded overflow-x-auto mb-2 text-left">
                    <pre class="text-xs whitespace-pre-wrap"><code>{sanitizeCodeSnippet(issue.extra.lines)}</code></pre>
                  </div>
                {/if} -->
                
                <!-- Location details -->
                <div class="mt-2 p-3 bg-base-200 rounded-lg">
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
                    <div>
                      <span class="text-xs font-bold text-primary-focus">File:</span>
                      <span class="text-xs text-base-content font-medium ml-1">{sanitizeText(issue.path)}</span>
                    </div>
                    <div>
                      <span class="text-xs font-bold text-primary-focus">Location:</span>
                      <span class="text-xs text-base-content font-medium ml-1">
                        Line {issue.start?.line || 0}
                        {#if issue.start?.col && issue.end?.col}
                          , Col {issue.start.col} - {issue.end.col}
                        {/if}
                      </span>
                    </div>
                  </div>
                  
                  <!-- Additional metadata -->
                  {#if issue.extra?.metadata}
                    <div class="mt-2 border-t border-base-300 pt-2">
                      {#if issue.extra.metadata.cwe}
                        <div class="mb-1">
                          <span class="text-xs font-bold text-primary-focus">CWE:</span>
                          <span class="text-xs text-base-content ml-1">
                            {Array.isArray(issue.extra.metadata.cwe) 
                              ? issue.extra.metadata.cwe.join(', ') 
                              : issue.extra.metadata.cwe}
                          </span>
                        </div>
                      {/if}
                      
                      {#if issue.extra.metadata.owasp}
                        <div class="mb-1">
                          <span class="text-xs font-bold text-primary-focus">OWASP:</span>
                          <span class="text-xs text-base-content ml-1">
                            {Array.isArray(issue.extra.metadata.owasp) 
                              ? issue.extra.metadata.owasp.join(', ') 
                              : issue.extra.metadata.owasp}
                          </span>
                        </div>
                      {/if}
                      
                      {#if issue.extra.metadata.references && 
                            (Array.isArray(issue.extra.metadata.references) && issue.extra.metadata.references.length > 0)}
                        <div class="mb-1">
                          <span class="text-xs font-bold text-primary-focus">References:</span>
                          <div class="ml-1">
                            {#each issue.extra.metadata.references as ref}
                              <a href={ref} target="_blank" class="text-xs link link-primary block truncate">{ref}</a>
                            {/each}
                          </div>
                        </div>
                      {/if}
                    </div>
                  {/if}

                </div>
              </div>
            </div>
          {/each}        
        </div>
        
        <!-- No results message for filtered view -->
        {#if filteredVulnerabilities.length === 0 && selectedSeverity !== 'ALL'}
          <div class="alert alert-info shadow-lg mt-4">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <div>
              <h3 class="font-bold">No vulnerabilities found</h3>
              <div class="text-sm">No vulnerabilities match the selected severity level: {selectedSeverity.toLowerCase()}</div>
            </div>
          </div>
        {/if}
      {/if}

      <!-- Parse Errors List -->
      {#if scan.report && scan.report.errors && scan.report.errors.length > 0}
        <div class="divider my-8"></div>
        <h4 class="text-lg font-bold mb-4 text-warning">Parsing Errors ({scan.report.errors.length})</h4>
        
        <div class="space-y-6">
          {#each scan.report.errors as error}
            <div class="card border border-warning shadow-md bg-base-100">
              <div class="card-body p-4">
                <div class="flex justify-between items-start flex-wrap gap-2">
                  <h5 class="card-title text-base text-warning">
                    Error Code: {error.code}
                  </h5>
                  <span class="badge badge-warning text-warning-content font-medium">
                    Level: {error.level}
                  </span>
                </div>
                
                <!-- Error message -->
                <div class="my-2 text-base-content font-medium">
                  {error.message || 'No description available'}
                </div>
                
                <!-- Error details -->
                <div class="mt-2 p-3 bg-base-200 rounded-lg">
                  <div class="grid grid-cols-1 gap-2">
                    {#if error.path}
                      <div>
                        <span class="text-xs font-bold text-warning">File:</span>
                        <span class="text-xs text-base-content font-medium ml-1">{error.path}</span>
                      </div>
                    {/if}
                    
                    {#if error.type && error.type.length > 0 && error.type[0]}
                      <div>
                        <span class="text-xs font-bold text-warning">Type:</span>
                        <span class="text-xs text-base-content font-medium ml-1">{error.type[0]}</span>
                      </div>
                    {/if}
                    
                    {#if error.spans && error.spans.length > 0}
                      <div class="mt-2 border-t border-base-300 pt-2">
                        <span class="text-xs font-bold text-warning block mb-1">Error Locations:</span>
                        <div class="ml-2 space-y-2">
                          {#each error.spans as span}
                            <div class="bg-base-300 p-2 rounded text-xs text-base-content">
                              <div><span class="font-semibold text-base-content">File:</span> {span.file}</div>
                              <div>
                                <span class="font-semibold text-base-content">Location:</span> 
                                Line {span.start.line}, Col {span.start.col} to 
                                Line {span.end.line}, Col {span.end.col}
                              </div>
                            </div>
                          {/each}
                        </div>
                      </div>
                    {/if}
                  </div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      <!-- Scanned Files -->
      {#if scan.report && scan.report.paths && scan.report.paths.scanned}
        <div class="mt-6">
          <h4 class="text-lg font-bold mb-2 text-primary-focus">Scanned Files</h4>
          <div class="bg-base-200 p-3 rounded-lg max-h-48 overflow-y-auto">
            <ul class="list-disc list-inside text-sm">
              {#each scan.report.paths.scanned as path}
                <li class="text-base-content">{path}</li>
              {/each}
            </ul>
          </div>
        </div>
      {/if}      {#if !scan.report?.results?.length && !scan.report?.errors?.length}
        <div class="alert alert-success shadow-lg mt-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          <div>
            <h3 class="font-bold">All Clear!</h3>
            <div class="text-sm">No vulnerabilities or errors were found in this scan!</div>
          </div>
        </div>
      {/if}
    </div>
      <!-- Modal Footer -->
    <div class="flex items-center justify-end p-4 border-t border-gray-200">
      <div class="mr-auto text-xs text-base-content opacity-70">
        {#if scan.report && scan.report.time && scan.report.time.total_time}
          Analysis completed in {(scan.report.time.total_time).toFixed(2)}s
        {/if}
      </div>
      <button class="btn btn-primary" on:click={closeModal}>Close</button>
    </div>
  </div>
</div>
