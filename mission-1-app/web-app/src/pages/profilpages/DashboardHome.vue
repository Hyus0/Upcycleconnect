<template>
    <div class="layout-wrapper public-dashboard">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > TABLEAU DE BORD</p>
                <h1 class="hero-title1">Bonjour {{ prenom }} 👋</h1>
                <p class="classic-text">
                    Voici un résumé de votre activité sur UpcycleConnect
                </p>
            </div>
            <router-link to="/profil/createAnnonce" class="btn-main-action">+ Déposer une annonce</router-link>
        </header>

        <div class="stats-grid dashboard-stats">
            <div class="card card--score">
                <p class="tag-score">UPCYCLING SCORE</p>
                <div class="score-value">{{ stats.total_points || 0 }} <span>pts</span></div>
                <p class="score-level">Niveau : {{ stats.niveau || 'Novice' }}</p>
                <div class="score-footer">
                    <div class="mini-stat">
                        <strong>{{ stats.co2_total_evite_kg || 0 }} kg</strong><br />CO2 évité
                    </div>
                    <div class="mini-stat">
                        <strong>{{ stats.nb_objets_recycles || 0 }}</strong><br />Objets
                    </div>
                    <div class="mini-stat">
                        <strong>EUR {{ stats.ressources_economisees || 0 }}</strong><br />Économisé
                    </div>
                </div>
            </div>
            
            <div class="card card--white">
                <div class="card-num">{{ annoncesActivesCount }}</div>
                <p class="text-dm">Annonces actives</p>
                <span class="badge badge--green">{{ annoncesActivesCount > 0 ? "+1 ce mois" : "Aucune active" }}</span>
            </div>

            <div class="card card--white">
                <div class="card-num2">{{ calendarEntries.length }}</div>
                <p class="text-dm">Inscriptions à venir</p>
                <span class="badge badge--orange">{{ nextEntryLabel }}</span>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Mes dernières annonces</h2>
                <div class="header-actions">
                    <router-link to="/profil/createAnnonce" class="btn-main-action1">+ Nouvelle annonce</router-link>
                </div>
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>OBJET</th>
                        <th>TYPE</th>
                        <th>STATUT PUBLICATION</th>
                        <th>STATUT LOGISTIQUE</th> 
                        <th>DATE</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="annonce in annonces.slice(0, 4)" :key="annonce.id">
                        <td>
                            <strong>{{ annonce.titre }}</strong><br>
                            <small class="table-subtext">{{ annonce.type_materiau }}</small>
                        </td>
                        <td>
                            <span :class="annonce.type === 'Don' ? 'tag-don' : 'tag-vente'">
                                {{ annonce.type === 'Don' ? 'DON' : 'VENTE ' + annonce.prix + 'EUR' }}
                            </span>
                        </td>
                        <td>
                            <span :class="annonce.est_valide === 'Valide' ? 'status-valid' : 'status-pending'">
                                {{ annonce.est_valide === 'Valide' ? 'APPROUVÉE' : 'EN ATTENTE' }}
                            </span>
                        </td>
                        <td>
                            <span class="status-logistique">
                                {{ annonce.statut || 'En attente' }}
                            </span>
                        </td>
                        <td>{{ formatDate(annonce.date_creation) }}</td>
                        <td class="actions-cell">
                            <button class="btn-view" @click="goToAnnonce(annonce.id)">Voir</button>
                            <button 
                                    v-if="annonce.est_valide === 'En attente'" 
                                    class="btn-modify" 
                                    @click="goToModify(annonce.id)"
                            >
                                    Modifier
                            </button>
                            <button class="btn-remove" @click="removeAnnonce(annonce.id)">Retirer</button>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div v-if="annonces.length === 0" class="state-card">Vous n'avez pas encore déposé d'annonces.</div>
        </div>

        <div class="section-container planning-section">
            <div class="section-header">
                <div>
                    <h2>Mon planning ({{ planningWeekLabel }})</h2>
                    <p class="classic-text">Cliquez pour voir le calendrier complet.</p>
                </div>
                <button class="btn-secondary" @click.stop="goToFullPlanning">Vue mensuelle</button>
            </div>

            <div class="planning-week">
                <div class="planning-week__weekday" v-for="day in compactWeekDays" :key="day.key">
                    <span>{{ day.weekLabel }}</span>
                    <strong>{{ day.dayLabel }}</strong>
                </div>
                <div
                    v-for="day in compactWeekDays"
                    :key="`${day.key}-card`"
                    class="planning-week__day"
                    :class="{ 'is-today': day.isToday }"
                >
                    <div class="planning-week__head">{{ day.dayNumber }}</div>
                    <div class="planning-week__entries">
                        <div
                            v-for="entry in day.entries.slice(0, 2)"
                            :key="entry.id"
                            class="planning-week__entry"
                            :class="`planning-entry--${entry.kind}`"
                        >
                            <strong>{{ entry.title }}</strong>
                            <span>{{ entry.timeLabel }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="calendarEntries.length === 0" class="state-card" style="margin-top: 1rem;">
                Le planning se remplira dès que vous rejoindrez une formation ou un événement.
            </div>
        </div>

        <!-- TIPS ET NOTIFICATIONS MIS À JOUR -->
        <div class="end-grid">
            <div class="section-container-tips">
                <p class="tag-conseil">💡 Conseil du jour</p>

                <div class="section-header">
                    <h2>Transformer un vieux jean en sac</h2>
                </div>
                <p>
                    Apprenez à confectionner un sac tote en 30 minutes avec un jean usé.
                    Matériel nécessaire : aiguille, fil, ciseaux.
                </p>
                <button class="btn-text-green">Lire la suite →</button>
            </div>
            <div class="section-container-tips">
                <p class="tag-notif">🔔 Notification</p>

                <div class="section-header">
                    <h2>Votre dépôt a été récupéré !</h2>
                </div>
                <p>
                    La chaise vintage que vous avez déposée le 12 fév. a été récupérée par un artisan.
                    +50 points Upcycling Score !
                </p>
                <button class="btn-text-green">Voir le projet →</button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const API_URL = "http://localhost:8081";

const prenom = ref(localStorage.getItem("userPrenom") || "Invité");
const annonces = ref([]);
const calendarEntries = ref([]);

const stats = ref({
    total_points: 0,
    niveau: "Chargement...",
    co2_total_evite_kg: 0,
    nb_objets_recycles: 0,
    ressources_economisees: 0,
});

const getLocalISODate = (date) => {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
};

const annoncesActivesCount = computed(() =>
    annonces.value.filter((a) => a.est_valide === "Valide" && a.statut === "Disponible").length
);

const nextEntryLabel = computed(() => {
    if (calendarEntries.value.length === 0) return "À venir";
    return [...calendarEntries.value].sort((a, b) => a.date.localeCompare(b.date))[0].dateLabel;
});

const currentWeekStart = computed(() => {
    const now = new Date();
    const offset = (now.getDay() + 6) % 7;
    const monday = new Date(now);
    monday.setHours(0, 0, 0, 0);
    monday.setDate(now.getDate() - offset);
    return monday;
});

const planningWeekLabel = computed(() =>
    currentWeekStart.value.toLocaleDateString("fr-FR", { day: "numeric", month: "long" })
);

const weekDays = ["Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];

const compactWeekDays = computed(() =>
    Array.from({ length: 7 }, (_, index) => {
        const date = new Date(currentWeekStart.value);
        date.setDate(currentWeekStart.value.getDate() + index);
        const isoDate = getLocalISODate(date);
        
        return {
            key: isoDate,
            weekLabel: weekDays[index].toUpperCase(),
            dayLabel: date.toLocaleDateString("fr-FR", { day: "numeric" }),
            dayNumber: date.getDate(),
            isToday: isoDate === getLocalISODate(new Date()),
            entries: calendarEntries.value.filter((e) => e.date === isoDate)
        };
    })
);

const formatDate = (val) => {
    if (!val) return "NULL";
    const date = new Date(val);
    return isNaN(date.getTime()) ? "NULL" : new Intl.DateTimeFormat("fr-FR", { day: "2-digit", month: "short", year: "numeric" }).format(date);
};

const goToAnnonce = (id) => router.push({ name: "see-annonce", params: { id } });
const goToFullPlanning = () => router.push("/profil/planning");

const removeAnnonce = async (id) => {
    if (!confirm("Voulez-vous vraiment retirer cette annonce ?")) return;
    try {
        const token = localStorage.getItem("userToken") || "";
        const res = await fetch(`${API_URL}/annonces/${id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } });
        if (!res.ok) throw new Error("Erreur suppression");
        annonces.value = annonces.value.filter((a) => a.id !== id);
    } catch (e) {
        alert("Suppression impossible pour le moment.");
    }
};

const loadCalendarEntries = async (id) => {
    try {
        const token = localStorage.getItem("userToken") || "";
        const res = await fetch(`${API_URL}/user/planning/${id}`, { headers: { 'Authorization': `Bearer ${token}` } });
        if (!res.ok) return;
        
        const data = await res.json();
        let entries = [];
        
        if (data && !Array.isArray(data)) {
            const mapEntry = (arr, kind, dateKey) => (arr || []).map(item => {
                const parsed = new Date(item[dateKey]);
                if (isNaN(parsed.getTime())) return null;
                return {
                    id: `${kind}-${item.id}`,
                    kind,
                    title: item.titre,
                    date: getLocalISODate(parsed),
                    dateLabel: parsed.toLocaleDateString("fr-FR", { day: "numeric", month: "short" }),
                    timeLabel: parsed.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })
                };
            }).filter(Boolean);
            
            entries = [...mapEntry(data.formations, 'formation', 'date_debut'), ...mapEntry(data.evenements, 'event', 'date_evenement')];
        }
        calendarEntries.value = entries.sort((a, b) => a.date.localeCompare(b.date));
    } catch (e) {
        console.error("Erreur planning", e);
    }
};

onMounted(async () => {
    const id = localStorage.getItem("userId");
    if (!id) return;
    const headers = { 'Authorization': `Bearer ${localStorage.getItem("userToken") || ""}` };

    try {
        const resStats = await fetch(`${API_URL}/users/${id}/stats`, { headers });
        if (resStats.ok) {
            const data = await resStats.json();
            stats.value = data;
            if (data.total_points !== undefined) {
                localStorage.setItem("userScore", String(data.total_points));
                window.dispatchEvent(new Event("auth-change"));
            }
        }
    } catch (e) {}

    try {
        const resAnnonces = await fetch(`${API_URL}/users/${id}/annonces`, { headers }); 
        if (resAnnonces.ok) {
            const data = await resAnnonces.json();
            annonces.value = Array.isArray(data) ? data : [];
        }
    } catch (e) {}

    await loadCalendarEntries(id);
});
</script>

<style scoped>
.public-dashboard {
    min-height: 100vh;
    padding: 20px;
    background: var(--bg-light, #f7f9f7);
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
}

.btn-main-action, .btn-main-action1 {
    display: inline-flex;
    align-items: center;
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border-radius: 10px;
    text-decoration: none;
    font-weight: bold;
    border: none;
    cursor: pointer;
}

.section-container {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 2rem;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
}

.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    color: var(--text-grey, #666);
    background: #fbfdfb;
    text-align: center;
}

.stats-grid {
    display: grid;
    gap: 1.5rem;
    margin-bottom: 2rem;
}

.dashboard-stats {
    grid-template-columns: 1.4fr 0.8fr 0.8fr; 
}

@media (max-width: 920px) {
    .dashboard-stats {
        grid-template-columns: 1fr;
    }
}

.card {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0,0,0,0.02);
}

.card--score {
    background: #2d7a4f;
    color: white;
    border: none;
}

.tag-score {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
}

.score-value {
    font-size: 2.8rem;
    font-weight: 900;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.score-value span {
    font-size: 1.2rem;
    font-weight: normal;
    opacity: 0.8;
}

.score-level {
    font-size: 1.1rem;
    margin-bottom: 1.5rem;
    opacity: 0.9;
}

.score-footer {
    display: flex;
    justify-content: space-between;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    padding-top: 1rem;
    margin-top: 1rem;
}

.mini-stat strong {
    font-size: 1.1rem;
    color: white;
    display: block;
}

.card-num, .card-num2 {
    font-size: 2.5rem;
    font-weight: bold;
    color: #1a1a1a;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.badge {
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: bold;
    display: inline-block;
}

.btn-remove {
    background: #ffe5e5;
    border: 1px solid #ffcccc;
    color: #d32f2f;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    margin-left: 10px;
    transition: all 0.2s;
}

.btn-remove:hover {
    background: #ffcccc; 
    color: #b71c1c;
}

.status-logistique {
    background: #f0f0f0;
    color: #555;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.badge--green { background: #eaf4ed; color: #2d7a4f; }
.badge--orange { background: #fff3cd; color: #856404; }

.table-subtext {
    display: block;
    margin-top: 4px;
    color: var(--text-grey, #666);
    font-size: 0.78rem;
}

.btn-view {
    background: transparent;
    border: 1px solid #ddd;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
}

.btn-view:hover {
    background: #f0f0f0;
}

.btn-remove {
    background: #FEE2E2;
    border: none;
    color: #dc3545;
    cursor: pointer;
    margin-left: 10px;
}

.btn-remove:hover {
    text-decoration: none;
}

.planning-week {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 1rem;
}

.planning-week__weekday {
    text-align: center;
    font-weight: 700;
    color: #666;
    padding-bottom: 5px;
    display: flex;
    flex-direction: column;
}

.planning-week__day {
    background: #fbfdfb;
    border: 1px solid #e0e0e0;
    border-radius: 12px;
    min-height: 100px;
    padding: 10px;
    display: flex;
    flex-direction: column;
}

.planning-week__day.is-today {
    border: 2px solid #2d7a4f;
    background: #eaf4ed;
}

.planning-week__head {
    font-weight: bold;
    font-size: 1.1rem;
    margin-bottom: 8px;
    color: #1a1a1a;
}

.planning-week__entry {
    padding: 6px;
    border-radius: 6px;
    font-size: 0.75rem;
    margin-bottom: 4px;
    display: flex;
    flex-direction: column;
}

.planning-entry--formation { background: #d1e7dd; color: #0f5132; }
.planning-entry--event { background: #fff3cd; color: #856404; }

.end-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
}

.section-container-tips {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
}

.tag-conseil {
    background: #eaf4ed;
    color: #2d7a4f;
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
}

.tag-notif {
    background: #FEF3C7;
    color: #92400E; 
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
}

.btn-text-green {
    background: #eaf4ed;
    border: none;
    color: #2d7a4f;
    font-weight: bold;
    font-size: 1rem;
    padding: 12px 20px;
    border-radius: 10px;
    margin-top: 1.5rem; 
    cursor: pointer;
    transition: all 0.2s ease;
    display: block; 
    text-align: center; 
}

.btn-text-green:hover {
    background: #d1e7dd; 
    color: #1b4d31;
    transform: translateY(-2px); 
    box-shadow: 0 4px 8px rgba(45, 122, 79, 0.15);
}
</style>