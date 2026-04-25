<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > DEPOTS</p>
            <h1 class="hero-title1">VOS DEPOTS CONTENEURS</h1>
            <p class="classic-text">
                Gérez vos réservations de casiers et accédez à vos codes de dépôt.
            </p>
        </div>
        <button class="btn-main-action" @click="$router.push('/profil/createAnnonce')">
            + Déposer une annonce
        </button>
    </header>

    <div v-if="loading" class="loading-state">Chargement...</div>

    <div v-else class="section-container">
        
        <div class="depot-group">
            <h2 class="group-title">📦 Prêts pour le dépôt <span class="count-badge">{{ aPlanifier.length }}</span></h2>
            <div v-if="aPlanifier.length === 0" class="empty-msg">Aucun objet en attente de planification.</div>
            
            <div class="annonces-grid">
                <div v-for="item in aPlanifier" :key="item.id" class="annonce-card">
                    <div class="card-header">
                        <span :class="item.type === 'Don' ? 'tag-don' : 'tag-vente'">
                            {{ item.type === 'Don' ? 'DON' : 'Vente ' + item.prix + '€' }}
                        </span>
                    </div>
                    <div class="card-body">
                        <h3 class="item-title">{{ item.titre }}</h3>
                        <p class="item-loc">📍 {{ item.ville }}</p>
                    </div>
                    <div class="card-footer">
                        <button class="btn-plan" @click="planifier(item.id)">
                            Réserver un casier
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div class="depot-group active-reservations">
            <h2 class="group-title">🔑 Réservations actives <span class="count-badge">{{ actifs.length }}</span></h2>
            <div v-if="actifs.length === 0" class="empty-msg">Aucune réservation en cours.</div>
            
            <div class="annonces-grid">
                <div v-for="depot in actifs" :key="depot.id" class="annonce-card reserved-border">
                    <div class="card-header">
                        <span class="tag-locker">CASIER N°{{ depot.id_casier }}</span>
                        <button class="btn-cancel" @click="annuler(depot.id)">Annuler</button>
                    </div>
                    <div class="card-body">
                        <h3 class="item-title">{{ depot.titre }}</h3>
                        <div class="pin-box">
                            <label>CODE PIN</label>
                            <div class="pin-code-small">{{ depot.code_pin_depot }}</div>
                        </div>
                    </div>
                    <div class="card-footer-info">
                        <p>📍 {{ depot.adresse }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router'; 

const router = useRouter();
const annonces = ref([]);
const loading = ref(true);

const aPlanifier = computed(() => 
    annonces.value.filter(a => a.est_valide === 'Valide' && a.statut === 'Disponible')
);

const actifs = computed(() => 
    annonces.value.filter(a => a.statut === 'En attente de depot')
);

const fetchAnnonces = async () => {
    const userId = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");
    try {
        const res = await fetch(`http://localhost:8081/users/${userId}/annonces`, {
            headers: { "Authorization": token }
        });
        if (res.ok) annonces.value = await res.json();
    } finally {
        loading.value = false;
    }
};

const planifier = (id) => { 
    router.push({ name: 'reserve-casier', params: { id: id } }); 
};

const annuler = async (id) => { 
    console.log("Libérer le casier pour l'annonce", id); 
};

onMounted(fetchAnnonces);
</script>

<style scoped>
.annonces-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
    margin-top: 1rem;
}

.annonce-card {
    background: white;
    border: 1px solid #eee;
    border-radius: 16px;
    padding: 1.2rem;
    display: flex;
    flex-direction: column;
    transition: transform 0.2s, box-shadow 0.2s;
}

.annonce-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 10px 20px rgba(0,0,0,0.05);
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.item-title {
    font-size: 1.1rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: #333;
}

.item-loc {
    font-size: 0.85rem;
    color: #888;
}

/* ÉLÉMENTS SPÉCIFIQUES DÉPÔT */
.depot-group { margin-bottom: 3rem; }

.group-title {
    font-size: 0.85rem;
    text-transform: uppercase;
    color: #aaa;
    letter-spacing: 1px;
    display: flex;
    align-items: center;
    gap: 10px;
}

.count-badge {
    background: #eee;
    color: #666;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 0.7rem;
}

.tag-locker {
    background: #1a1c23;
    color: #ffd43b;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.7rem;
    font-weight: 900;
}

.pin-box {
    background: #f8f9fa;
    border-radius: 10px;
    padding: 10px;
    text-align: center;
    margin-top: 1rem;
    border: 1px dashed #ddd;
}

.pin-box label {
    font-size: 0.6rem;
    font-weight: 800;
    color: #999;
}

.pin-code-small {
    font-size: 1.5rem;
    font-weight: 900;
    letter-spacing: 4px;
    color: #2d6a4f;
}

/* BOUTONS */
.btn-plan {
    width: 100%;
    background: #2d7a4f;
    color: white;
    border: none;
    padding: 10px;
    border-radius: 8px;
    font-weight: 700;
    cursor: pointer;
    margin-top: 1rem;
}

.btn-cancel {
    background: none;
    border: none;
    color: #e63946;
    font-size: 0.75rem;
    font-weight: 700;
    cursor: pointer;
}

.card-footer-info {
    margin-top: auto;
    padding-top: 10px;
    font-size: 0.75rem;
    color: #999;
}

.reserved-border {
    border: 2px solid #1a1c23;
}

.empty-msg {
    padding: 2rem;
    color: #bbb;
    font-style: italic;
    font-size: 0.9rem;
}

.tag-don { background: #e8f5e9; color: #2e7d32; font-weight: 800; padding: 4px 10px; border-radius: 6px; font-size: 0.7rem; }
.tag-vente { background: #fff8e1; color: #f57f17; font-weight: 800; padding: 4px 10px; border-radius: 6px; font-size: 0.7rem; }
</style>