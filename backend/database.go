package main

func getTier(tierURL string) string {
	row := db.QueryRow(`SELECT tier FROM tiers WHERE url = ?`, tierURL)
	var tier string

	if err := row.Scan(&tier); err != nil {
		return ""
	}

	return tier
}
