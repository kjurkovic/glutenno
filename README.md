# Glutenno Coolinarika

Aplikacija sluzi za prikaz recepata i uputa za kuhanje ljudima oboljelima od celijakije (alergije na gluten) ili one koji se tako osjecaju :)


## Funkcionalnosti
1. REST API za recepte
2. Autorizacija i registracija korisnika
3. REST API za komentare registriranih korisnika
4. Servis za slanje mailova

Svi backend servisi se spajaju na jednu PostgreSQL bazu podataka ali svaki servis ima odvojenu schemu. Kroz config.yml pojedinacnog servisa moguce je konfigurirati odvojenu bazu podataka.
Konfiguracija za svaki pojedinacni backend servis se nalazi u datoteci {imeservisa}/config.yml

## Pokretanje aplikacije

Za pokrenuti aplikaciju potrebno je imati instaliran i pokrenut Docker.
Aplikacija se pokrece preko development.sh skripte koja se nalazi u `scripts` direktoriju:

```sh
$ ./scripts/development.sh
```

Skripta ce pokrenuti make nad svim servisima koji ce kreirati docker image svakog pojedinacnog servisa. Nakon toga ce se pozvati `docker compose up` cime se dizu servisi definirani u `./scripts/docker-compose.yml` datoteci kao i postgresql baza podataka.