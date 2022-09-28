<script context="module">
    let current_id = 0;
</script>

<script>
    let id = ++current_id;

    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    function switchTab() {
        dispatch('switch', {
            tab: id
        });
    }
</script>

<style>
.tab {
  float: left;
}
.tab-switch {
  display: none;
}
</style>

<div class="tab">
    {#if id === 1}
    <input type="radio" name="css-tabs" id="tab-{id}" checked class="tab-switch" on:change={switchTab}>
    {:else}
    <input type="radio" name="css-tabs" id="tab-{id}" class="tab-switch" on:change={switchTab}>
    {/if}
    <slot>
    </slot>
</div>
