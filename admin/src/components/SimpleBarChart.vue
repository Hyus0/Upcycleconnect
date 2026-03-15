<template>
  <article class="surface-card chart-card">
    <header>
      <h3>{{ title }}</h3>
      <p>{{ subtitle }}</p>
    </header>

    <div class="bars">
      <div v-for="item in items" :key="item.label" class="bar-row">
        <div class="label-row">
          <span>{{ item.label }}</span>
          <strong>{{ item.missing ? "API" : item.value }}</strong>
        </div>
        <div class="track">
          <div
            class="fill"
            :class="{ missing: item.missing }"
            :style="{ width: `${getWidth(item)}%` }"
          ></div>
        </div>
      </div>
    </div>
  </article>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    required: true
  },
  subtitle: {
    type: String,
    default: ""
  },
  items: {
    type: Array,
    default: () => []
  }
});

function getWidth(item) {
  const max = Math.max(...props.items.map((entry) => entry.value || 0), 1);
  if (item.missing) {
    return 15;
  }
  return Math.max(8, Math.round(((item.value || 0) / max) * 100));
}
</script>

<style scoped>
.chart-card {
  padding: 20px;
}

header h3 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

header p {
  margin: 8px 0 0;
  color: var(--text-secondary);
}

.bars {
  display: grid;
  gap: 16px;
  margin-top: 18px;
}

.label-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.track {
  height: 14px;
  border-radius: 999px;
  background: rgba(45, 122, 79, 0.08);
  overflow: hidden;
}

.fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, var(--brand-green), var(--accent-teal));
}

.fill.missing {
  background: linear-gradient(90deg, var(--accent-amber), var(--accent-coral));
}
</style>
