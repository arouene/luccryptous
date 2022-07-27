<script context="module">
	let current_id = 0;
</script>

<script>
    let id = ++current_id;
</script>

<style>
.tab-label {
  position: relative;
  display: block;
  line-height: 2.75em;
  height: 3em;
  padding: 0 1.618em;
  background: #424242;
  border-right: 0.125rem solid #919191;
  color: #fff;
  cursor: pointer;
  top: 0;
  transition: all 0.25s;
}
.tab-label:hover {
  top: -0.25rem;
  transition: top 0.25s;
}
:global(.tab-switch:checked + .tab-label) {
  background: #fff;
  color: #2c3e50;
  border-bottom: 0;
  border-right: 0.125rem solid #fff;
  transition: all 0.35s;
  z-index: 1;
  top: -0.0625rem;
}
</style>

<label for="tab-{id}" class="tab-label">
	<slot>
	</slot>
</label>
