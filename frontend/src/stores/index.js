import { writable } from 'svelte/store';

// Local storage persistence for stores
function createPersistentStore(key, initialValue) {
  // Create a writable store
  const store = writable(initialValue);
  
  // Check for browser environment (needed for SSR compatibility)
  const isBrowser = typeof window !== 'undefined';
  
  if (isBrowser) {
    // Check if we have stored data
    const storedValue = localStorage.getItem(key);
    
    if (storedValue) {
      // If we have stored data, use it for the initial value
      store.set(JSON.parse(storedValue));
    }
    
    // Subscribe to changes and update localStorage
    store.subscribe(value => {
      localStorage.setItem(key, JSON.stringify(value));
    });
  }
  
  return store;
}

// Sidebar state
export const sidebarOpen = writable(true);

// Current active menu
export const activeMenu = writable('scan'); // 'scan' or 'history'

// Scanning state
export const isScanning = writable(false);

// History of scans - persisted to localStorage
export const scanHistory = createPersistentStore('scanHistory', []);
// export const scanHistory = writable([])

