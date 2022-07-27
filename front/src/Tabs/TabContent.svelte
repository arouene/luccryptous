<script>
    export let height = undefined;
</script>

<style>
.tab-content {
  position: absolute;
  z-index: 1;
  top: 2.75em;
  left: 0;
  width: 100%;
  padding: 1.618rem;
  background: #fff;
  color: #2c3e50;
  border-bottom: 0.25rem solid #bdc3c7;
  opacity: 0;
  transition: all 0.35s;
}
:global(.tab-switch:checked + label + .tab-content) {
  z-index: 2;
  opacity: 1;
  transition: all 0.35s;
}
</style>

<div class="tab-content" style="height: {height}">
  <slot>
  </slot>
</div>
