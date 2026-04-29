export const fallbackAnnonces = [
  {
    id: 1,
    titre: "Chaise vintage en bois",
    description: "Chaise a restaurer avec structure solide.",
    statut: "en ligne",
    prix: 35,
    etat_objet: "bon etat",
    ville: "Paris",
    adresse: "11e arrondissement",
    code_postal: "75011",
    type: "vente",
    date_creation: "2026-03-12T10:00:00Z"
  },
  {
    id: 2,
    titre: "Lot de bocaux en verre",
    description: "Bocaux propres, parfaits pour atelier zero dechet.",
    statut: "en ligne",
    prix: 0,
    etat_objet: "tres bon etat",
    ville: "Lyon",
    adresse: "Lyon 2",
    code_postal: "69002",
    type: "don",
    date_creation: "2026-03-11T10:00:00Z"
  }
];

export { fetchAnnonces } from "./publicApi";
