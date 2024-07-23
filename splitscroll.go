package go_blueberry

const SplitScroll string = `
{{ define "SplitScroll" }}
<template id="split-scroll-template" class="split-scroll">
  <style>
    .split-scroll {
      display: flex;
      overflow: hidden;
      max-height: 100vh;
    }
  </style>
    <div class="split-scroll">
      <slot name="stay"> </slot>
      <slot name="scroll"> </slot>
    </div>
</template>

<script>
  customElements.define(
    "split-scroll",
    class extends HTMLElement {
      constructor() {
        super();
        const template = document.getElementById("split-scroll-template").content;
        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(template.cloneNode(true));
      }

      connectedCallback() {
        const stay = this.querySelector("[slot='stay']");
        const scroll = this.querySelector("[slot='scroll']");
        scroll.style.overflow = "scroll";
        scroll.style.scrollbarWidth = "none";
        stay.addEventListener("wheel", (event) => {
          let curPos = scroll.scrollTop;
          let nextPos = curPos - event.wheelDeltaY;
          scroll.scroll(0, nextPos)
          if (nextPos > 0 && nextPos < scroll.scrollTopMax) {
            event.preventDefault();
          }
        });
      }
    },
  );
</script>
{{ end }}
`
