import Dexie from "dexie"

const DB_NAME = "offlineWalletDB"

const db = new Dexie(DB_NAME)

db.version(1).stores({ table_wallet: '++id, name, mainCoinId, chainType, epkey, address, isFromSystem, uuid' })

export default db;
