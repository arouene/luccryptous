<style>
:global(*) {
  box-sizing: border-box;
}
.wrapper {
  max-width: 50rem;
  width: 100%;
  margin: 0 auto;
}
.tabs {
  position: relative;
  margin: 3rem 0;
  background: transparent;
}
.tabs::before,
.tabs::after {
  content: "";
  display: table;
}
.tabs::after {
  clear: both;
}
</style>

<div class="wrapper">
  <div class="tabs">
    <slot></slot>
  </div>
</div>
