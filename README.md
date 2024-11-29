# Test

The tests expect `ECONOMIC_AGREEMENT_GRANT_TOKEN` and
`ECONOMIC_APP_SECRET_TOKEN` to be set in the environment variables. Be /very/
aware that the test is not read-only; it will try to create stuff, so either
run it only on a test account, or look through every test to verify it does no
harm.

I (aks) don't know if it makes a lot of sense to do pure testing of this
library?
