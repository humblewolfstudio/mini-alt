<template>
  <div class="drive-card">
    <div class="gauge-container">
      <svg viewBox="0 0 36 36" class="circular-chart">
        <path
            class="circle-bg"
            d="M18 2.0845
             a 15.9155 15.9155 0 0 1 0 31.831
             a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <path
            class="circle"
            :stroke-dasharray="usagePercentDisplay + ', 100'"
            d="M18 2.0845
             a 15.9155 15.9155 0 0 1 0 31.831
             a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <text x="18" y="17.35" class="percentage">
          {{ usagePercentDisplay }}%
        </text>
        <text x="18" y="22" class="label">Used Capacity</text>
      </svg>
    </div>

    <div class="drive-info">
      <div class="capacity-row">
        <div>
          <div class="capacity-label">Used Capacity</div>
          <div class="capacity-value">
            {{ formatSize(specs?.UsedCapacity) }}
          </div>
        </div>

        <div class="divider"></div>

        <div>
          <div class="capacity-label">Available Capacity</div>
          <div class="capacity-value">
            {{ formatSize(specs?.FreeCapacity) }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { SystemSpecsResponse } from "../sources/SystemDataSource";
import {formatSize} from "../utils";

const props = defineProps<{
  specs: SystemSpecsResponse | null
}>();

const usagePercentDisplay = computed(() => {
  if (!props.specs) return 0;
  return ((props.specs.UsedCapacity / props.specs.TotalCapacity) * 100).toFixed(2);
});

const freePercentDisplay = computed(() => {
  if (!props.specs) return 0;
  return ((props.specs.FreeCapacity / props.specs.TotalCapacity) * 100).toFixed(2);
});
</script>

<style scoped>
.drive-card {
  display: flex;
  align-items: center;
  flex-direction: row;
  gap: 30px;
  padding: 20px;
  background: white;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  border-radius: 8px;
}

.gauge-container {
  flex-shrink: 0;
}

.circular-chart {
  display: block;
  width: 140px;
  height: 140px;
}

.circle-bg {
  fill: none;
  stroke: #2c3e50;
  stroke-width: 3.8;
}

.circle {
  fill: none;
  stroke-width: 3.8;
  stroke-linecap: round;
  stroke: #42b983;
  transition: stroke-dasharray 0.6s ease;
}

.percentage {
  font-size: 0.25em;
  text-anchor: middle;
  font-weight: bold;
}

.label {
  font-size: 0.2em;
  text-anchor: middle;
}

.drive-info {
  flex: 1;
}

.drive-info h3 {
  margin: 0 0 10px;
  font-size: 1.1rem;
}

.capacity-row {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.capacity-label {
  font-size: 0.9rem;
}

.capacity-value {
  color: #42b983;
  font-size: 1.5rem;
  font-weight: bold;
}

.capacity-sub {
  font-size: 0.85rem;
}

.divider {
  width: 1px;
  background: #334155;
  align-self: stretch;
}

@media (max-width: 600px) {
  .drive-card {
    flex-direction: column;
  }
}
</style>
