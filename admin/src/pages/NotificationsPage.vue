<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Notifications</div>
        <h2 class="page-title">Campagnes et messages</h2>
        <p class="page-description">Preparation, planification et suivi des notifications.</p>
      </div>
      <button class="button button-secondary" @click="loadNotifications">Actualiser</button>
    </header>

    <div class="split-grid">
      <article class="surface-card section-card stack">
        <div class="panel-head">
          <h3>Nouvelle notification</h3>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>
        <FormField label="Titre">
          <input v-model="form.title" placeholder="Nouveau depot recupere" />
        </FormField>
        <div class="filters-grid">
          <FormField label="Canal">
            <BaseSelect v-model="form.channel" :options="channelOptions" />
          </FormField>
          <FormField label="Audience">
            <BaseSelect v-model="form.audience" :options="audienceOptions" />
          </FormField>
          <FormField label="Statut">
            <BaseSelect v-model="form.status" :options="statusOptions.slice(1)" />
          </FormField>
        </div>
        <FormField label="Planification">
          <input v-model="form.scheduledAt" type="datetime-local" />
        </FormField>
        <FormField label="Message">
          <textarea v-model="form.message"></textarea>
        </FormField>
        <div class="toolbar">
          <button class="button button-primary" @click="submitForm">Creer</button>
        </div>
      </article>

      <article class="surface-card section-card stack">
        <div class="panel-head">
          <h3>Filtres</h3>
          <span class="mini-note">{{ pagination?.total ?? 0 }} notifications</span>
        </div>
        <div class="filters-grid">
          <FormField label="Recherche">
            <input v-model="filters.search" placeholder="Titre ou message" />
          </FormField>
          <FormField label="Canal">
            <BaseSelect v-model="filters.channel" :options="filterChannelOptions" />
          </FormField>
          <FormField label="Statut">
            <BaseSelect v-model="filters.status" :options="statusOptions" />
          </FormField>
        </div>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadNotifications" />
    <EmptyState v-else-if="rows.length === 0" title="Aucune notification" message="Aucune campagne ne correspond aux filtres." />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-title="{ row }">
        <div class="identity">
          <strong>{{ row.title }}</strong>
          <span>{{ row.message }}</span>
        </div>
      </template>
      <template #cell-status="{ row }">
        <StatusBadge :label="row.status" :tone="row.status === 'sent' ? 'green' : 'amber'" />
      </template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="setStatus(row, 'scheduled')">Programmer</button>
          <button class="button button-primary" @click="setStatus(row, 'sent')">Envoyer</button>
          <button class="button button-danger" @click="remove(row.id)">Supprimer</button>
        </div>
      </template>
    </DataTable>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import BaseSelect from "../components/BaseSelect.vue";
import DataTable from "../components/DataTable.vue";
import EmptyState from "../components/EmptyState.vue";
import ErrorState from "../components/ErrorState.vue";
import FormField from "../components/FormField.vue";
import LoadingState from "../components/LoadingState.vue";
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const filters = reactive({ search: "", channel: "", status: "", page: 1, pageSize: 8 });
const form = reactive({ title: "", channel: "email", audience: "all", status: "draft", scheduledAt: "", message: "" });

const columns = [
  { key: "title", label: "Notification" },
  { key: "channel", label: "Canal" },
  { key: "audience", label: "Audience" },
  { key: "status", label: "Statut" },
  { key: "scheduledAt", label: "Planifiee le" }
];

const channelOptions = [
  { label: "Email", value: "email" },
  { label: "Push", value: "push" },
  { label: "SMS", value: "sms" }
];

const filterChannelOptions = computed(() => [{ label: "Tous", value: "" }, ...channelOptions]);

const audienceOptions = [
  { label: "Tous", value: "all" },
  { label: "Particuliers", value: "particuliers" },
  { label: "Prestataires", value: "prestataires" },
  { label: "Admins", value: "admins" }
];

const statusOptions = [
  { label: "Tous", value: "" },
  { label: "Brouillon", value: "draft" },
  { label: "Programmee", value: "scheduled" },
  { label: "Envoyee", value: "sent" }
];

function resetForm() {
  form.title = "";
  form.channel = "email";
  form.audience = "all";
  form.status = "draft";
  form.scheduledAt = "";
  form.message = "";
}

async function loadNotifications() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listNotifications(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les notifications.";
  } finally {
    loading.value = false;
  }
}

async function submitForm() {
  try {
    await adminApi.createNotification({ ...form });
    pushToast({ title: "Notification creee", message: "La campagne a ete ajoutee.", tone: "green" });
    resetForm();
    await loadNotifications();
  } catch (err) {
    pushToast({ title: "Creation impossible", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function setStatus(row, status) {
  try {
    await adminApi.updateNotificationStatus(row.id, status);
    pushToast({ title: "Statut mis a jour", message: `${row.title} a ete mise a jour.`, tone: "green" });
    await loadNotifications();
  } catch (err) {
    pushToast({ title: "Mise a jour impossible", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function remove(id) {
  try {
    await adminApi.deleteNotification(id);
    pushToast({ title: "Notification supprimee", message: "Suppression effectuee.", tone: "coral" });
    await loadNotifications();
  } catch (err) {
    pushToast({ title: "Suppression impossible", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function changePage(page) {
  filters.page = page;
}

watch(() => [filters.search, filters.channel, filters.status, filters.page], loadNotifications);
onMounted(loadNotifications);
</script>

<style scoped>
.identity {
  display: grid;
  gap: 4px;
}

.identity span,
.mini-note {
  color: var(--text-secondary);
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.panel-head h3 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

.actions {
  justify-content: flex-end;
}
</style>
