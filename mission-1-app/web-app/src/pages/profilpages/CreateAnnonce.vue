<template>
  <header class="content-header">
    <div class="header-left">
      <p class="sidebar__category2">ACCUEIL > ANNONCES > CRÉER</p>
      <h1 class="hero-title1">DÉPOSER UN OBJET</h1>
      <p class="classic-text">Remplissez les détails pour donner une seconde vie à votre objet.</p>
    </div>
    <div v-if="errors.length > 0" class="error-box">
        <ul style="margin: 0; padding-left: 20px;">
            <li v-for="(err, index) in errors" :key="index">{{ err }}</li>
        </ul>
    </div>
    <div v-if="successMsg" class="success-box">
        {{ successMsg }}
    </div>
    <button class="btn-secondary" @click="$router.back()">🠔 Retour</button>
  </header>

  <form @submit.prevent="handleSubmit" class="create-annonce-form">
    <div class="split-layout">
      
      <div class="form-card main-info-card">
        <h2 class="card-title">1. L'Objet</h2>
        <p class="card-subtitle">Décrivez précisément ce que vous souhaitez céder.</p>

        <div class="form-group">
          <label>Titre de l'annonce</label>
          <input v-model="form.titre" type="text" required placeholder="Ex: Bureau en chêne massif" />
        </div>

        <div class="form-group">
          <label>Description</label>
          <textarea v-model="form.description" rows="6" placeholder="Décrivez l'état, les dimensions, les éventuels défauts..."></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Catégorie</label>
            <select v-model="form.id_categorie" required>
              <option value="" disabled>Choisir une catégorie</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.nom }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>État de l'objet</label>
            <select v-model="form.etat_objet" required>
              <option value="Neuf">Neuf</option>
              <option value="Bon etat">Bon état</option>
              <option value="Usage">À restaurer / Usagé</option>
            </select>
          </div>
        </div>

        <div class="form-divider">Transaction</div>

        <div class="form-row align-center">
          <div class="form-group">
            <label>Type d'annonce</label>
            <div class="radio-group">
              <label class="radio-label">
                <input type="radio" v-model="form.type" value="Don" /> Don
              </label>
              <label class="radio-label">
                <input type="radio" v-model="form.type" value="Vente" /> Vente
              </label>
            </div>
          </div>

          <div class="form-group price-group" v-if="form.type === 'Vente'">
            <label>Prix (€)</label>
            <div class="price-input-wrapper">
              <input v-model.number="form.prix" type="number" step="0.01" min="0" placeholder="0.00" />
              <span class="currency-symbol">€</span>
            </div>
          </div>
        </div>
      </div>

      <div class="right-column">
        
        <div class="form-card tech-info-card">
          <h2 class="card-title">2. Détails Upcycling</h2>
          <p class="card-subtitle">Ces infos nous aident à calculer votre impact environnemental.</p>

          <div class="form-row">
            <div class="form-group">
              <label>Matériau principal</label>
              <input v-model="form.type_materiau" type="text" placeholder="Ex: Bois, Métal, Plastique" />
            </div>
            <div class="form-group">
              <label>Poids estimé (kg)</label>
              <input v-model.number="form.poids_estime_kg" type="number" step="0.1" placeholder="0.0" />
            </div>
          </div>
        </div>

        <div class="form-card geo-info-card">
          <h2 class="card-title">3. Localisation</h2>
          <p class="card-subtitle">Où l'objet doit-il être déposé ?</p>

          <div class="form-row">
            <div class="form-group">
              <label>Ville</label>
              <input v-model="form.ville" type="text" required placeholder="Paris" />
            </div>
            <div class="form-group">
              <label>Code Postal</label>
              <input v-model="form.code_postal" type="text" maxlength="5" required placeholder="75001" />
            </div>
          </div>

          <div class="form-group">
            <label>Adresse précise</label>
            <input v-model="form.adresse" type="text" placeholder="Rue, numéro... pour le site de dépôt" />
          </div>
        </div>

        <div class="form-actions-card">
          <button type="button" class="btn-cancel" @click="$router.back()">Annuler</button>
          <button type="submit" class="btn-save" :disabled="loading">
            {{ loading ? 'Publication...' : 'Publier l\'annonce ✓' }}
          </button>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const loading = ref(false);
const categories = ref([]);  

const errors = ref([]);
const successMsg = ref("");

const form = ref({
  id_vendeur: parseInt(sessionStorage.getItem("userId")) || 0,
  id_categorie: "",
  titre: "",
  description: "",
  type_materiau: "",
  poids_estime_kg: 0,
  prix: 0,
  etat_objet: "Bon etat",
  type: "Don",
  ville: "",
  code_postal: "",
  adresse: "",
  statut: "Disponible"
});

onMounted(async () => {
  try {
    const res = await fetch("http://localhost:8081/categories");
    if (res.ok) categories.value = await res.json();
  } catch (err) {
    console.error("Erreur categories:", err);
  }
});

const validateFrontend = () => {
  errors.value = [];
  if (form.value.titre.length < 5) errors.value.push("Le titre doit faire au moins 5 caractères.");
  if (!form.value.id_categorie) errors.value.push("Veuillez choisir une catégorie.");
  if (form.value.type === 'Vente' && form.value.prix <= 0) errors.value.push("Le prix doit être supérieur à 0€.");
  return errors.value.length === 0;
};

const handleSubmit = async () => {
    loading.value = true;
    errors.value = [];
    successMsg.value = "";

    const token = sessionStorage.getItem("userToken");

    try {
        const response = await fetch("http://localhost:8081/annonces", {
            method: "POST",
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(form.value)
        });

        if (response.ok) {
            successMsg.value = "Annonce créée avec succès !";
            setTimeout(() => {
                router.push("/profil/annonces");
            }, 2000);
        } else {
            const errorData = await response.json();
            
            if (errorData.errors) {
                errors.value = errorData.errors;
            } else {
                alert("Erreur : " + (errorData.message || "La création a échoué."));
            }
        }
    } catch (err) {
        console.error("Erreur lors de l'envoi:", err);
        alert("Serveur injoignable.");
    } finally {
        loading.value = false;
    }
};
</script>

<style scoped>

.error-box {
    background-color: #fee2e2;
    border: 1px solid #ef4444;
    color: #b91c1c;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
}

.success-box {
    background-color: #f0fdf4;
    border: 1px solid #22c55e;
    color: #166534;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
}

.create-annonce-form {
  width: 100%;
  max-width: none;
  box-sizing: border-box;
}

.split-layout {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 2rem;
  width: 100%;
}

.form-card {
  background: #ffffff;
  padding: 2rem;
  border-radius: 16px;
  border: 1px solid #eee;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.right-column {
  display: flex;
  flex-direction: column;
}

.form-actions-card {
  margin-top: auto; 
  padding: 1rem 0;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.card-title {
  font-size: 1.4rem;
  font-weight: 700;
  margin: 0;
  color: #1a1a1a;
}

.card-subtitle {
  font-size: 0.9rem;
  color: #666;
  margin: -1rem 0 1rem 0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.form-divider {
  margin-top: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--brand-mist);
  font-weight: 800;
  color: var(--brand-green);
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 1px;
}

input, select, textarea {
  padding: 0.9rem;
  border: 1px solid #ddd;
  border-radius: 10px;
  font-family: inherit;
  font-size: 0.95rem;
  transition: border-color 0.2s;
}

input:focus, select:focus, textarea:focus {
  outline: none;
  border-color: var(--brand-green);
}

textarea {
  resize: vertical;
}

.radio-group {
  display: flex;
  gap: 1.5rem;
  padding: 0.5rem 0;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  cursor: pointer;
  font-weight: 600;
}

.radio-label input[type="radio"] {
  width: 1.2rem;
  height: 1.2rem;
  accent-color: var(--brand-green);
}

.price-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.price-input-wrapper input {
  width: 100%;
  padding-right: 2.5rem;
}

.currency-symbol {
  position: absolute;
  right: 1rem;
  font-weight: bold;
  color: #888;
}

.btn-save {
  background-color: var(--brand-green);
  color: white;
  padding: 1rem 2.5rem;
  border: none;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-save:hover {
  background-color: #246343;
}

.btn-save:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.btn-cancel {
  background: none;
  border: 1px solid #ccc;
  color: #666;
  padding: 1rem 2.5rem;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
}

.btn-cancel:hover {
  background-color: #f5f5f5;
  border-color: #999;
}
</style>