<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Gestion des evenements</div>
        <h2 class="page-title">Planning admin en attente d'API dediee</h2>
        <p class="page-description">
          Les filtres et composants sont prets, mais il manque encore la ressource evenements pour
          permettre la recherche par date et le CRUD complet.
        </p>
      </div>
    </header>

    <div class="split-grid">
      <article class="surface-card section-card stack">
        <div class="filters-grid">
          <FormField label="Recherche">
            <input disabled placeholder="Titre ou lieu" />
          </FormField>
          <FormField label="Date">
            <DatePickerField model-value="" />
          </FormField>
          <FormField label="Statut">
            <BaseSelect :model-value="''" :options="statusOptions" />
          </FormField>
        </div>
      </article>

      <ResourceNotice
        title="Endpoints evenements a exposer"
        message="Ces routes permettront la liste, la recherche par date et le CRUD demande dans le cahier des charges."
        :items="requiredEndpoints"
      />
    </div>

    <EmptyState
      title="Aucun evenement disponible"
      message="Le backend n'expose pas encore GET /api/events. La page reste volontairement honnete sur cet etat."
    />
  </section>
</template>

<script setup>
import BaseSelect from "../components/BaseSelect.vue";
import DatePickerField from "../components/DatePickerField.vue";
import EmptyState from "../components/EmptyState.vue";
import FormField from "../components/FormField.vue";
import ResourceNotice from "../components/ResourceNotice.vue";

const statusOptions = [
  { label: "Tous les statuts", value: "" },
  { label: "Planifie", value: "planned" },
  { label: "Publie", value: "published" },
  { label: "Archive", value: "archived" }
];

const requiredEndpoints = [
  "GET /api/events",
  "GET /api/events/:id",
  "POST /api/events",
  "PUT /api/events/:id",
  "DELETE /api/events/:id"
];
</script>
