<template>
  <section class="stack">
    <header class="page-header">
      <div>
        <div class="eyebrow">Gestion des categories</div>
        <h2 class="page-title">Taxonomie prete, API encore absente</h2>
        <p class="page-description">
          La page admin est structuree pour accueillir une hierarchie parent/enfant, mais aucun
          endpoint categories n'est expose aujourd'hui dans le repo courant.
        </p>
      </div>
    </header>

    <div class="split-grid">
      <ResourceNotice
        title="Endpoints categories a ajouter"
        message="Pour rendre la page pleinement operationnelle, voici le minimum necessaire cote backend."
        :items="requiredEndpoints"
      />

      <article class="surface-card section-card stack">
        <h3>Formulaire prevu</h3>
        <FormField label="Nom de la categorie">
          <input placeholder="Ex: Mobilier exterieur" disabled />
        </FormField>
        <FormField label="Categorie parente">
          <BaseSelect :model-value="''" :options="parentOptions" />
        </FormField>
        <FormField label="Description">
          <textarea disabled placeholder="Description et usage de la categorie"></textarea>
        </FormField>
        <button class="button button-primary" disabled>Enregistrer</button>
      </article>
    </div>

    <EmptyState
      title="Aucune categorie chargee"
      message="Le composant de listing est pret, mais l'API n'expose pas encore GET /api/categories."
    />
  </section>
</template>

<script setup>
import BaseSelect from "../components/BaseSelect.vue";
import EmptyState from "../components/EmptyState.vue";
import FormField from "../components/FormField.vue";
import ResourceNotice from "../components/ResourceNotice.vue";

const requiredEndpoints = [
  "GET /api/categories",
  "GET /api/categories/:id",
  "POST /api/categories",
  "PUT /api/categories/:id",
  "DELETE /api/categories/:id"
];

const parentOptions = [{ label: "Aucune categorie parente", value: "" }];
</script>
