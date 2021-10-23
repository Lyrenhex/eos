/**
 * StorageSpace class
 * 
 * This class handles interfacing with whichever storage provider the
 * application is currently using, and can be swapped out solely by
 * replacing this implementation with another (but must provide the
 * same interfaces).
 * 
 * StorageSpace implements:
 * - setItem(key, value): Store a string with name 'key'.
 * - getItem(key): Get the string with name 'key', or return `null`.
 * - clear(): Deletes all data stored.
 * - capacity(): Returns the capacity of the storage space in bytes.
 * - usage(): Returns the number of bytes used from the capacity in bytes.
 */

class StorageSpace {
    #storage;
    #usage;
    #capacity;
    #usageChangeCallback;

    /**
     * Constructs a new Storage interface, handling setting up a new
     * (or loading an existing) storage space.
     * @param {function} usageChangeCallback a function called whenever the
     * usage of the storage space increases or decreases; this function must
     * accept a single parameter of type number.
     */
    constructor(usageChangeCallback) {
        this.#storage = window.localStorage;
        this.#usage = this.usage();
        this.#usageChangeCallback = usageChangeCallback;
    }

    /**
     * Calculates the length of a key-value pair in bytes.
     * @param {string} key 
     * @param {string} value 
     * @returns {number} the length of the new key-value pair in bytes.
     */
    #size(key, value) {
        return ((key + value) * 2) / 1024;
    }

    /**
     * setItem stores the string 'value' with name 'key'.
     * 
     * Note: errors not caused by quota being exceeded are
     * silently logged, and the rest of the application will
     * not be aware anything went wrong. The user will instead
     * be alerted outside of the application flow.
     * 
     * @param {string} key 
     * @param {string} value 
     * @returns {boolean} true if the value was saved, or
     * false if the storage space was full.
     */
    setItem(key, value) {
        let deltaSize = this.#size(key, value);
        if (this.getItem(key) !== null) {
            deltaSize = (value.length - this.getItem(key).length) * 2;
        }
        try {
            this.#storage.setItem(key, value);
        } catch (err) {
            if (err.code == 22 || err.name === "NS_ERROR_DOM_QUOTA_REACHED") {
                // quota full
                console.warn("LocalStorage quota full. Cleanup routine required.");
                return false;
            }
            console.error(err);
            alert(`Unexpected error ${err.code} ${err.name} ${err.description}. Data modifications have NOT been saved. Please raise this as an issue on GitHub (you are also advised to back up your Eos data).`);
            return true;
        }
        this.#usage += deltaSize;
        this.#usageChangeCallback(this.#usage);
        return true;
    }

    /**
     * getItem gets the string value stored with name 'key'.
     * @param {string} key 
     * @returns {string} the value stored under name 'key'.
     */
    getItem(key) {
        return this.#storage.getItem(key);
    }

    /**
     * clear deletes all key-value pairs stored in the storage space.
     */
    clear() {
        this.#storage.clear();
        this.#storage.setItem('capacity', this.capacity());
        this.#usage = null;
        this.#usage = this.usage();
        console.log("calling callback", this.#usage);
        this.#usageChangeCallback(this.#usage);
        console.log("called");
    }

    /**
     * capacity calculates the maximum size of the storage space.
     * @returns {number} the maximum size of the storage space in bytes.
     */
    capacity() {
        if (!this.#capacity) {
            this.#capacity = this.getItem('capacity');
            if (!this.#capacity) {
                var i = 0;
                try {
                    // Test up to 10 MB
                    for (i = 250; i <= 10000; i += 250) {
                        this.#storage.setItem('test', new Array((i * 1024) + 1).join('a'));
                    }
                } catch (e) {
                    this.#storage.removeItem('test');
                    this.#capacity = ((i - 250) + Math.round(this.#usage / 1024 / 250) * 250) * 1024;
                    this.#storage.setItem('capacity', this.#capacity);
                }
            }
        }

        return this.#capacity;
    }

    /**
     * usage calculates the current byte usage of the storage space.
     * @returns {number} the number of bytes currently used in the storage space.
     */
    usage() {
        if (!this.#usage) {
            var allStrings = '';
            for (var key in this.#storage) {
                if (this.#storage.hasOwnProperty(key)) {
                    allStrings += key + this.#storage[key];
                }
            }
            this.#usage = allStrings ? allStrings.length * 2 : 0;
        }

        console.log(this.#usage);

        return this.#usage;
    }
}