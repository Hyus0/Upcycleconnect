<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Ressources</div>
        <h2 class="page-title">Fonctionnalites front</h2>
        <p class="page-description">
          Administration CRUD des donnees qui alimentent les parcours publics, profils, logistique, paiement,
          messagerie, forum, abonnements et traductions.
        </p>
      </div>
      <div class="toolbar">
        <button class="button button-secondary" @click="loadRows">Actualiser</button>
      </div>
    </header>

    <div class="split-grid admin-workspace">
      <article class="surface-card section-card stack admin-panel admin-panel--form">
        <div class="panel-head">
          <div>
            <span class="panel-kicker">Ressource</span>
            <h3>{{ currentResource?.label || "Selection" }}</h3>
          </div>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>

        <FormField label="Module gere">
          <BaseSelect v-model="selectedResourceKey" :options="resourceOptions" />
        </FormField>

        <div v-if="currentResource" class="mode-note">
          <strong>{{ editingKey ? "Edition" : "Creation" }}</strong>
          <span>{{ writableFields.length }} champs modifiables sur cette ressource.</span>
        </div>

        <div v-if="currentResource" class="dynamic-form">
          <FormField v-for="field in writableFields" :key="field.name" :label="field.label">
            <BaseSelect
              v-if="field.type === 'select'"
              v-model="form[field.name]"
              :options="selectOptions(field)"
            />
            <textarea
              v-else-if="field.type === 'textarea'"
              v-model="form[field.name]"
              :required="field.required"
            ></textarea>
            <select v-else-if="field.type === 'boolean'" v-model="form[field.name]">
              <option value="1">Oui</option>
              <option value="0">Non</option>
            </select>
            <input
              v-else
              v-model="form[field.name]"
              :type="inputType(field)"
              :required="field.required"
            />
          </FormField>
        </div>

        <div class="toolbar">
          <button class="button button-primary" :disabled="!currentResource || saving" @click="submitForm">
            {{ editingKey ? "Enregistrer" : "Creer" }}
          </button>
        </div>
      </article>

      <article class="surface-card section-card stack admin-panel admin-panel--filters">
        <div class="panel-head panel-head--compact">
          <div>
            <span class="panel-kicker">Couverture</span>
            <h3>Modules administrables</h3>
          </div>
          <span class="mini-note">{{ resources.length }} ressources</span>
        </div>
        <div class="admin-stats">
          <div class="admin-stat">
            <strong>{{ rows.length }}</strong>
            <span>Lignes chargees</span>
          </div>
          <div class="admin-stat">
            <strong>{{ currentResource?.fields?.length || 0 }}</strong>
            <span>Champs</span>
          </div>
          <div class="admin-stat">
            <strong>{{ currentResource?.primaryKeys?.length || 0 }}</strong>
            <span>Cles</span>
          </div>
        </div>
        <FormField label="Filtre rapide">
          <input v-model="search" placeholder="Rechercher dans les lignes chargees" />
        </FormField>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadRows" />
    <EmptyState v-else-if="filteredRows.length === 0" title="Aucune donnee" message="Aucune ligne pour cette ressource." />
    <DataTable v-else row-key="__key" :columns="columns" :rows="pagedRows" :pagination="pagination" @page-change="page = $event">
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="startEdit(row)">Editer</button>
          <button class="button button-danger" @click="rowToDelete = row">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Supprimer cette ligne ?"
      :message="deleteMessage"
      confirm-label="Supprimer"
      @cancel="rowToDelete = null"
      @confirm="deleteCurrent"
    />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import BaseSelect from "../components/BaseSelect.vue";
import ConfirmModal from "../components/ConfirmModal.vue";
import DataTable from "../components/DataTable.vue";
import EmptyState from "../components/EmptyState.vue";
import ErrorState from "../components/ErrorState.vue";
import FormField from "../components/FormField.vue";
import LoadingState from "../components/LoadingState.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";

const { pushToast } = useToastStore();

const resources = ref([]);
const selectedResourceKey = ref("");
const rows = ref([]);
const loading = ref(false);
const saving = ref(false);
const error = ref("");
const editingKey = ref("");
const rowToDelete = ref(null);
const search = ref("");
const page = ref(1);
const pageSize = 10;
const form = reactive({});

const currentResource = computed(() =>
  resources.value.find((item) => item.key === selectedResourceKey.value)
);

const resourceOptions = computed(() =>
  resources.value.map((item) => ({ label: item.label, value: item.key }))
);

const writableFields = computed(() =>
  (currentResource.value?.fields || []).filter((field) => !field.readOnly)
);

const visibleFields = computed(() => (currentResource.value?.fields || []).slice(0, 6));

const columns = computed(() =>
  visibleFields.value.map((field) => ({ key: field.name, label: field.label }))
);

const filteredRows = computed(() => {
  const needle = search.value.trim().toLowerCase();
  if (!needle) return rows.value;
  return rows.value.filter((row) =>
    Object.values(row).join(" ").toLowerCase().includes(needle)
  );
});

const pagedRows = computed(() => {
  const start = (page.value - 1) * pageSize;
  return filteredRows.value.slice(start, start + pageSize);
});

const pagination = computed(() => ({
  page: page.value,
  pageSize,
  total: filteredRows.value.length,
  totalPages: Math.max(1, Math.ceil(filteredRows.value.length / pageSize))
}));

const deleteMessage = computed(() => {
  if (!rowToDelete.value || !currentResource.value) return "";
  return `${currentResource.value.label} - ${rowToDelete.value.__key}`;
});

watch(selectedResourceKey, async () => {
  resetForm();
  page.value = 1;
  search.value = "";
  await loadRows();
});

watch(filteredRows, () => {
  if (page.value > pagination.value.totalPages) {
    page.value = pagination.value.totalPages;
  }
});

onMounted(async () => {
  await loadResources();
});

async function loadResources() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listResources();
    resources.value = (response.items || []).sort((a, b) => a.label.localeCompare(b.label));
    selectedResourceKey.value = resources.value[0]?.key || "";
    if (selectedResourceKey.value) {
      await loadRows();
    }
  } catch (err) {
    error.value = err.message || "Impossible de charger les ressources.";
  } finally {
    loading.value = false;
  }
}

async function loadRows() {
  if (!selectedResourceKey.value) return;
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listResourceRows(selectedResourceKey.value);
    rows.value = response.items || [];
  } catch (err) {
    error.value = err.message || "Impossible de charger la ressource.";
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  editingKey.value = "";
  Object.keys(form).forEach((key) => delete form[key]);
  writableFields.value.forEach((field) => {
    form[field.name] = field.type === "boolean" ? "0" : "";
  });
}

function startEdit(row) {
  editingKey.value = row.__key;
  Object.keys(form).forEach((key) => delete form[key]);
  writableFields.value.forEach((field) => {
    form[field.name] = row[field.name] ?? (field.type === "boolean" ? "0" : "");
  });
}

async function submitForm() {
  if (!currentResource.value) return;
  saving.value = true;
  try {
    const payload = {};
    writableFields.value.forEach((field) => {
      payload[field.name] = form[field.name];
    });
    if (editingKey.value) {
      await adminApi.updateResourceRow(selectedResourceKey.value, editingKey.value, payload);
      pushToast({ title: "Ligne mise a jour", message: "Modification enregistree.", tone: "green" });
    } else {
      await adminApi.createResourceRow(selectedResourceKey.value, payload);
      pushToast({ title: "Ligne creee", message: "Nouvelle donnee ajoutee.", tone: "green" });
    }
    resetForm();
    await loadRows();
  } catch (err) {
    pushToast({ title: "Operation impossible", message: err.message || "Operation impossible.", tone: "coral" });
  } finally {
    saving.value = false;
  }
}

async function deleteCurrent() {
  if (!rowToDelete.value) return;
  try {
    await adminApi.deleteResourceRow(selectedResourceKey.value, rowToDelete.value.__key);
    pushToast({ title: "Ligne supprimee", message: "Suppression effectuee.", tone: "coral" });
    rowToDelete.value = null;
    await loadRows();
  } catch (err) {
    pushToast({ title: "Suppression impossible", message: err.message || "Suppression impossible.", tone: "coral" });
  }
}

function selectOptions(field) {
  return [
    { label: "Non renseigne", value: "" },
    ...(field.options || []).map((option) => ({ label: option, value: option }))
  ];
}

function inputType(field) {
  if (field.type === "number" || field.type === "decimal") return "number";
  if (field.type === "datetime") return "datetime-local";
  return "text";
}
</script>

<style scoped>
.dynamic-form {
  display: grid;
  gap: 14px;
}
</style>
