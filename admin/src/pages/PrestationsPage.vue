<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Catalogue</div>
        <h2 class="page-title">Annonces / Prestations</h2>
        <p class="page-description">Creation et gestion des objets, dons, ventes et prestations via l'API.</p>
      </div>
      <div class="toolbar">
        <button class="button button-primary" @click="showCreate = !showCreate">
          {{ showCreate ? "Fermer" : "Nouvelle ligne" }}
        </button>
        <button class="button button-secondary" @click="loadPrestations">Actualiser</button>
      </div>
    </header>

    <div class="split-grid catalogue-workspace">
      <article class="surface-card section-card stack catalogue-panel catalogue-panel--filters">
        <div class="panel-head panel-head--compact">
          <div>
            <span class="panel-kicker">Pilotage</span>
            <h3>Catalogue actif</h3>
          </div>
          <span class="mini-note">{{ pagination?.total ?? 0 }} lignes</span>
        </div>

        <div class="catalogue-stats">
          <div class="catalogue-stat">
            <strong>{{ pagination?.total ?? rows.length }}</strong>
            <span>Total</span>
          </div>
          <div class="catalogue-stat">
            <strong>{{ publishedCount }}</strong>
            <span>Publiees</span>
          </div>
          <div class="catalogue-stat">
            <strong>{{ draftCount }}</strong>
            <span>Brouillons</span>
          </div>
        </div>

        <div class="type-shortcuts" aria-label="Filtres rapides">
          <button
            v-for="option in typeFilterOptions"
            :key="option.value || 'all'"
            class="type-chip"
            :class="{ active: filters.type === option.value }"
            type="button"
            @click="filters.type = option.value"
          >
            {{ option.label }}
          </button>
        </div>

        <div class="filters-compact">
          <FormField label="Recherche">
            <input v-model="filters.search" placeholder="Titre ou description" />
          </FormField>
          <FormField label="Tri">
            <BaseSelect v-model="filters.sortBy" :options="sortOptions" />
          </FormField>
        </div>

        <button class="button button-secondary filters-reset" type="button" @click="resetFilters">
          Reinitialiser les filtres
        </button>
      </article>

      <article v-if="showCreate" class="surface-card section-card stack catalogue-panel catalogue-panel--form">
        <div class="panel-head">
          <div>
            <span class="panel-kicker">{{ form.kind === "annonce" ? "Objet" : "Service" }}</span>
            <h3>{{ editingId ? "Modifier" : "Ajouter" }}</h3>
          </div>
          <button class="button button-ghost" @click="resetForm">Vider</button>
        </div>

        <div class="mode-note">
          <strong>{{ form.kind === "annonce" ? "Annonce objet" : "Prestation catalogue" }}</strong>
          <span>
            {{
              form.kind === "annonce"
                ? "Creation d'un don ou d'une vente rattachee a un vendeur."
                : "Creation d'un service, d'une vente catalogue ou d'un don publie par l'admin."
            }}
          </span>
        </div>

        <div class="filters-grid form-grid">
          <FormField label="Mode de creation">
            <BaseSelect v-model="form.kind" :options="kindOptions" />
          </FormField>
          <FormField label="Titre" :error="formErrors.title">
            <input v-model="form.title" />
          </FormField>
          <FormField v-if="form.kind === 'prestation'" label="Prestataire" :error="formErrors.provider">
            <input v-model="form.provider" placeholder="Atelier Renouveau" />
          </FormField>
          <FormField v-else label="ID vendeur" :error="formErrors.sellerId">
            <input v-model="form.sellerId" type="number" min="1" />
          </FormField>
          <FormField label="Type">
            <BaseSelect v-model="form.type" :options="currentTypeOptions" />
          </FormField>
          <FormField label="Prix" :error="formErrors.price" :hint="priceHint">
            <input v-model="form.price" type="number" min="0" step="0.01" />
          </FormField>
          <FormField label="Statut">
            <BaseSelect v-model="form.status" :options="currentStatusOptions" />
          </FormField>
          <FormField v-if="form.kind === 'annonce'" label="Validation">
            <BaseSelect v-model="form.validation" :options="validationOptions" />
          </FormField>
          <FormField v-if="form.kind === 'annonce'" label="Etat de l'objet" :error="formErrors.condition">
            <BaseSelect v-model="form.condition" :options="conditionOptions" />
          </FormField>
          <FormField v-if="form.kind === 'annonce'" label="Ville" :error="formErrors.city">
            <input v-model="form.city" placeholder="Paris" />
          </FormField>
          <FormField v-if="form.kind === 'annonce'" label="Code postal" :error="formErrors.postalCode">
            <input v-model="form.postalCode" maxlength="5" placeholder="75011" />
          </FormField>
          <FormField v-if="form.kind === 'annonce'" label="Adresse" :error="formErrors.address">
            <input v-model="form.address" placeholder="12 rue de l'Upcycling" />
          </FormField>
        </div>
        <FormField label="Description" :error="formErrors.description">
          <textarea v-model="form.description"></textarea>
        </FormField>
        <div class="toolbar">
          <button class="button button-primary" @click="submitForm">
            {{ editingId ? "Enregistrer" : submitLabel }}
          </button>
        </div>
      </article>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" retry-label="Recharger" @retry="loadPrestations" />
    <EmptyState v-else-if="rows.length === 0" title="Aucune ligne" message="Aucun resultat." />
    <DataTable v-else :columns="columns" :rows="rows" :pagination="pagination" @page-change="changePage">
      <template #cell-title="{ row }">
        <div class="identity">
          <strong>{{ row.title }}</strong>
          <span>{{ row.description }}</span>
        </div>
      </template>
      <template #cell-price="{ row }">{{ formatCurrency(row.price) }}</template>
      <template #cell-status="{ row }"><StatusBadge :label="row.status" tone="teal" /></template>
      <template #actions="{ row }">
        <div class="toolbar actions">
          <button class="button button-secondary" @click="startEdit(row)">Editer</button>
          <button class="button button-danger" @click="confirmDelete(row)">Supprimer</button>
        </div>
      </template>
    </DataTable>

    <ConfirmModal
      :open="Boolean(rowToDelete)"
      title="Supprimer cette prestation ?"
      :message="rowToDelete ? rowToDelete.title : ''"
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
import StatusBadge from "../components/StatusBadge.vue";
import { adminApi } from "../services/api";
import { useToastStore } from "../store/toastStore";
import { formatCurrency } from "../utils/format";

const { pushToast } = useToastStore();

const loading = ref(true);
const error = ref("");
const rows = ref([]);
const pagination = ref(null);
const showCreate = ref(true);
const editingId = ref("");
const rowToDelete = ref(null);

const filters = reactive({ search: "", type: "", sortBy: "title", page: 1, pageSize: 7 });
const form = reactive({
  kind: "prestation",
  title: "",
  description: "",
  provider: "",
  sellerId: 1,
  type: "service",
  price: 0,
  status: "draft",
  validation: "En attente",
  condition: "Bon etat",
  address: "",
  city: "",
  postalCode: ""
});
const formErrors = reactive({
  title: "",
  description: "",
  price: "",
  provider: "",
  sellerId: "",
  condition: "",
  address: "",
  city: "",
  postalCode: ""
});

const columns = [
  { key: "title", label: "Prestation" },
  { key: "provider", label: "Prestataire" },
  { key: "type", label: "Type" },
  { key: "price", label: "Prix", align: "right" },
  { key: "status", label: "Statut" }
];

const kindOptions = [
  { label: "Prestation", value: "prestation" },
  { label: "Annonce objet", value: "annonce" }
];
const prestationTypeOptions = [
  { label: "Service", value: "service" },
  { label: "Vente", value: "vente" },
  { label: "Don", value: "don" }
];
const annonceTypeOptions = [
  { label: "Don", value: "Don" },
  { label: "Vente", value: "Vente" }
];
const typeFilterOptions = [{ label: "Tous", value: "" }, ...prestationTypeOptions];
const sortOptions = [
  { label: "Titre", value: "title" },
  { label: "Prix croissant", value: "price-asc" },
  { label: "Prix decroissant", value: "price-desc" }
];
const statusOptions = [
  { label: "Tous", value: "" },
  { label: "Brouillon", value: "draft" },
  { label: "Publie", value: "published" },
  { label: "Archive", value: "archived" }
];
const annonceStatusOptions = [
  { label: "Disponible", value: "Disponible" },
  { label: "Reserve", value: "Reserve" },
  { label: "Vendu", value: "Vendu" },
  { label: "Annule", value: "Annule" }
];
const validationOptions = [
  { label: "En attente", value: "En attente" },
  { label: "Valide", value: "Valide" },
  { label: "Refuse", value: "Refuse" }
];
const conditionOptions = [
  { label: "Neuf", value: "Neuf" },
  { label: "Tres bon etat", value: "Tres bon etat" },
  { label: "Bon etat", value: "Bon etat" },
  { label: "A reparer", value: "A reparer" }
];

const currentTypeOptions = computed(() => (form.kind === "annonce" ? annonceTypeOptions : prestationTypeOptions));
const currentStatusOptions = computed(() => (form.kind === "annonce" ? annonceStatusOptions : statusOptions.slice(1)));
const submitLabel = computed(() => (form.kind === "annonce" ? "Creer l'annonce" : "Creer la prestation"));
const publishedCount = computed(() => rows.value.filter((row) => row.status === "published" || row.status === "Disponible").length);
const draftCount = computed(() => rows.value.filter((row) => row.status === "draft" || row.status === "Brouillon").length);

const priceHint = computed(() => {
  if (form.kind === "annonce" && form.type === "Don") {
    return "Un don doit rester a 0 euro.";
  }
  if (form.kind === "prestation" && form.type === "don") {
    return "Un don doit rester a 0 euro.";
  }
  if (form.kind === "prestation" && form.status === "published") {
    return "Une prestation publiee doit avoir un prix superieur a 0.";
  }
  return "";
});

function validateForm() {
  formErrors.title = form.title.trim().length < 3 ? "Titre trop court." : "";
  formErrors.description = form.description.trim().length < 10 ? "Description trop courte." : "";
  formErrors.price = Number(form.price) < 0 ? "Prix invalide." : "";
  formErrors.provider = "";
  formErrors.sellerId = "";
  formErrors.condition = "";
  formErrors.address = "";
  formErrors.city = "";
  formErrors.postalCode = "";

  if (form.kind === "prestation" && form.status === "published" && form.provider.trim().length < 2) {
    formErrors.provider = "Le prestataire est obligatoire pour une publication.";
  }
  if ((form.kind === "prestation" && form.type === "don" || form.kind === "annonce" && form.type === "Don") && Number(form.price) !== 0) {
    formErrors.price = "Un don doit rester a 0 euro.";
  }
  if (form.kind === "prestation" && (form.type === "service" || form.type === "vente") && form.status === "published" && Number(form.price) <= 0) {
    formErrors.price = "Une prestation publiee doit avoir un prix superieur a 0.";
  }
  if (form.kind === "annonce") {
    formErrors.sellerId = Number(form.sellerId) <= 0 ? "Vendeur invalide." : "";
    formErrors.condition = form.condition ? "" : "Etat obligatoire.";
    formErrors.address = form.address.trim().length < 3 ? "Adresse obligatoire." : "";
    formErrors.city = form.city.trim().length < 2 ? "Ville obligatoire." : "";
    formErrors.postalCode = /^\d{5}$/.test(form.postalCode) ? "" : "Code postal invalide.";
  }
  return Object.values(formErrors).every((message) => !message);
}

function resetForm() {
  editingId.value = "";
  form.kind = "prestation";
  form.title = "";
  form.description = "";
  form.provider = "";
  form.sellerId = 1;
  form.type = "service";
  form.price = 0;
  form.status = "draft";
  form.validation = "En attente";
  form.condition = "Bon etat";
  form.address = "";
  form.city = "";
  form.postalCode = "";
  Object.keys(formErrors).forEach((key) => {
    formErrors[key] = "";
  });
}

function startEdit(row) {
  showCreate.value = true;
  editingId.value = row.id;
  form.kind = "prestation";
  form.title = row.title;
  form.description = row.description;
  form.provider = row.provider;
  form.type = row.type;
  form.price = row.price;
  form.status = row.status;
  window.scrollTo({ top: 0, behavior: "smooth" });
}

async function submitForm() {
  if (!validateForm()) return;
  try {
    if (form.kind === "annonce") {
      await adminApi.createAnnonce({ ...form });
      pushToast({ title: "Annonce creee", message: "Nouvelle annonce objet ajoutee.", tone: "green" });
    } else if (editingId.value) {
      await adminApi.updatePrestation(editingId.value, { ...form });
      pushToast({ title: "Prestation mise a jour", message: "Modification enregistree.", tone: "green" });
    } else {
      await adminApi.createPrestation({ ...form });
      pushToast({ title: "Prestation creee", message: "Nouvelle ligne ajoutee.", tone: "green" });
    }
    resetForm();
    await loadPrestations();
  } catch (err) {
    pushToast({ title: "Echec de l'enregistrement", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

async function loadPrestations() {
  loading.value = true;
  error.value = "";
  try {
    const response = await adminApi.listPrestations(filters);
    rows.value = response.items;
    pagination.value = response.pagination;
  } catch (err) {
    error.value = err.message ?? "Impossible de charger les prestations.";
  } finally {
    loading.value = false;
  }
}

function confirmDelete(row) {
  rowToDelete.value = row;
}

async function deleteCurrent() {
  if (!rowToDelete.value) return;
  try {
    await adminApi.deletePrestation(rowToDelete.value.id);
    pushToast({ title: "Prestation supprimee", message: "Suppression effectuee.", tone: "coral" });
    rowToDelete.value = null;
    await loadPrestations();
  } catch (err) {
    pushToast({ title: "Echec de suppression", message: err.message ?? "Operation impossible.", tone: "coral" });
  }
}

function changePage(page) {
  filters.page = page;
}

function resetFilters() {
  filters.search = "";
  filters.type = "";
  filters.sortBy = "title";
  filters.page = 1;
}

watch(() => [filters.search, filters.type, filters.sortBy, filters.page], loadPrestations);
watch(() => form.kind, (kind) => {
  if (kind === "annonce") {
    form.type = "Don";
    form.status = "Disponible";
    form.price = 0;
  } else {
    form.type = "service";
    form.status = "draft";
  }
});
watch(() => form.type, (type) => {
  if ((form.kind === "annonce" && type === "Don") || (form.kind === "prestation" && type === "don")) {
    form.price = 0;
  }
});

onMounted(loadPrestations);
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

.catalogue-workspace {
  grid-template-columns: minmax(320px, 0.82fr) minmax(520px, 1.18fr);
  align-items: start;
}

.catalogue-panel {
  min-height: auto;
  border-color: rgba(139, 210, 166, 0.16);
  background:
    radial-gradient(circle at 18% 0%, rgba(58, 158, 91, 0.14), transparent 34%),
    rgba(7, 22, 16, 0.88);
}

.catalogue-panel--filters {
  position: sticky;
  top: 24px;
  gap: 22px;
}

.catalogue-panel--form {
  background:
    linear-gradient(135deg, rgba(24, 63, 43, 0.34), transparent 42%),
    rgba(7, 22, 16, 0.92);
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.panel-head--compact {
  align-items: flex-start;
}

.panel-head h3 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

.panel-kicker {
  display: block;
  margin-bottom: 8px;
  color: #76e6a0;
  font-family: "Space Mono", monospace;
  font-size: 0.68rem;
  letter-spacing: 0.3em;
  text-transform: uppercase;
}

.catalogue-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.catalogue-stat {
  display: grid;
  gap: 6px;
  padding: 16px;
  border: 1px solid rgba(139, 210, 166, 0.16);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.035);
}

.catalogue-stat strong {
  color: var(--text-primary);
  font-family: "Syne", sans-serif;
  font-size: clamp(1.45rem, 3vw, 2.2rem);
  line-height: 1;
}

.catalogue-stat span {
  color: var(--text-secondary);
  font-size: 0.82rem;
}

.type-shortcuts {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.type-chip {
  border: 1px solid rgba(139, 210, 166, 0.18);
  border-radius: 999px;
  padding: 10px 14px;
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.035);
  cursor: pointer;
  transition: border-color 0.2s ease, color 0.2s ease, background 0.2s ease;
}

.type-chip:hover,
.type-chip.active {
  border-color: rgba(118, 230, 160, 0.5);
  color: var(--text-primary);
  background: rgba(46, 125, 82, 0.32);
}

.filters-compact {
  display: grid;
  gap: 16px;
}

.filters-reset {
  justify-self: start;
}

.mode-note {
  display: grid;
  gap: 6px;
  padding: 16px 18px;
  border: 1px solid rgba(139, 210, 166, 0.14);
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.035);
}

.mode-note strong {
  color: var(--text-primary);
}

.mode-note span {
  color: var(--text-secondary);
}

.form-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.actions {
  justify-content: flex-end;
}

@media (max-width: 1180px) {
  .catalogue-workspace,
  .form-grid {
    grid-template-columns: 1fr;
  }

  .catalogue-panel--filters {
    position: static;
  }
}

@media (max-width: 680px) {
  .catalogue-stats {
    grid-template-columns: 1fr;
  }
}
</style>
