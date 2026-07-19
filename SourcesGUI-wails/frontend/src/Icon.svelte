<script>
  export let name = '';
  export let size = 18;
</script>

<svg
  class="app-icon"
  width={size}
  height={size}
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="1.8"
  stroke-linecap="round"
  stroke-linejoin="round"
  aria-hidden="true"
>
  {#if name === 'sparkles'}
    <path d="M12 6c0 4-2 6-6 6 4 0 6 2 6 6 0-4 2-6 6-6-4 0-6-2-6-6Z" />
    <path d="M5.5 2.5c0 1.7-1.3 3-3 3 1.7 0 3 1.3 3 3 0-1.7 1.3-3 3-3-1.7 0-3-1.3-3-3Z" />
    <path d="M18.5 3.5c0 1.4-1.1 2.5-2.5 2.5 1.4 0 2.5 1.1 2.5 2.5 0-1.4 1.1-2.5 2.5-2.5-1.4 0-2.5-1.1-2.5-2.5Z" />
    <path d="M5 16c0 1.1-.9 2-2 2 1.1 0 2 .9 2 2 0-1.1.9-2 2-2-1.1 0-2-.9-2-2Z" />
  {:else if name === 'scripts'}
    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z" />
    <path d="M14 2v6h6M9 13l-2 2 2 2M15 13l2 2-2 2" />
  {:else if name === 'archive'}
    <path d="M4 7h16v13a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V7ZM3 2h18v5H3zM9 12h6" />
  {:else if name === 'font'}
    <path d="M4 7V4h16v3M9 20h6M12 4v16" />
  {:else if name === 'image'}
    <rect x="3" y="3" width="18" height="18" rx="2" />
    <circle cx="8.5" cy="8.5" r="1.5" />
    <path d="m21 15-5-5L5 21" />
  {:else if name === 'dialogues'}
    <path d="M21 15a4 4 0 0 1-4 4H8l-5 3V7a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4Z" />
    <path d="M8 9h8M8 13h5" />
  {:else if name === 'tools'}
    <path d="M14.7 6.3a4 4 0 0 0-5-5L12 3.6 8.4 7.2 6.1 4.9a4 4 0 0 0 5 5L19 17.8a2.1 2.1 0 0 1-3 3l-7.9-7.9a4 4 0 0 0-5-5l2.3 2.3-3.6 3.6-2.3-2.3a4 4 0 0 0 5 5" />
  {:else if name === 'info'}
    <circle cx="12" cy="12" r="9" />
    <path d="M12 11v5M12 8h.01" />
  {:else if name === 'terminal'}
    <rect x="3" y="4" width="18" height="16" rx="2" />
    <path d="m7 9 3 3-3 3M13 15h4" />
  {:else if name === 'chevron-down'}
    <path d="m6 9 6 6 6-6" />
  {:else if name === 'chevron-up'}
    <path d="m18 15-6-6-6 6" />
  {:else if name === 'trash'}
    <path d="M3 6h18M8 6V4h8v2M19 6l-1 15H6L5 6M10 11v5M14 11v5" />
  {:else if name === 'stop'}
    <rect x="6" y="6" width="12" height="12" rx="1" />
  {:else if name === 'check'}
    <circle cx="12" cy="12" r="9" />
    <path d="m8 12 2.7 2.7L16.5 9" />
  {:else if name === 'warning'}
    <path d="M10.3 3.6 2.4 17.4A2 2 0 0 0 4.1 20h15.8a2 2 0 0 0 1.7-2.6L13.7 3.6a2 2 0 0 0-3.4 0Z" />
    <path d="M12 9v4M12 17h.01" />
  {:else}
    <circle cx="12" cy="12" r="9" />
  {/if}
</svg>
