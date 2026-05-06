//// [tests/cases/compiler/commentsFormatting.ts] ////

//// [commentsFormatting.ts]
namespace m {
    /** this is first line - aligned to class declaration
* this is 4 spaces left aligned
 * this is 3 spaces left aligned
  * this is 2 spaces left aligned
   * this is 1 spaces left aligned
    * this is at same level as first line
     * this is 1 spaces right aligned
      * this is 2 spaces right aligned
       * this is 3 spaces right aligned
        * this is 4 spaces right aligned
         * this is 5 spaces right aligned
          * this is 6 spaces right aligned
           * this is 7 spaces right aligned
            * this is 8 spaces right aligned */
    export class c {
    }

        /** this is first line - 4 spaces right aligned to class but in js file should be aligned to class declaration
* this is 8 spaces left aligned
 * this is 7 spaces left aligned
  * this is 6 spaces left aligned
   * this is 5 spaces left aligned
    * this is 4 spaces left aligned
     * this is 3 spaces left aligned
      * this is 2 spaces left aligned
       * this is 1 spaces left aligned
        * this is at same level as first line
         * this is 1 spaces right aligned
          * this is 2 spaces right aligned
           * this is 3 spaces right aligned
            * this is 4 spaces right aligned
             * this is 5 spaces right aligned
              * this is 6 spaces right aligned
               * this is 7 spaces right aligned
                * this is 8 spaces right aligned */
    export class c2 {
    }

    /** this is comment with new lines in between

this is 4 spaces left aligned but above line is empty

 this is 3 spaces left aligned but above line is empty

  this is 2 spaces left aligned but above line is empty

    this is 1 spaces left aligned but above line is empty

    this is at same level as first line but above line is empty 

     this is 1 spaces right aligned but above line is empty

      this is 2 spaces right aligned but above line is empty

       this is 3 spaces right aligned but above line is empty

        this is 4 spaces right aligned but above line is empty
    
    
    Above 2 lines are empty
    
    
    
    above 3 lines are empty*/
    export class c3 {
    }

    /** this is first line - aligned to class declaration
	*              this is 0 space + tab
 	*              this is 1 space + tab
  	*              this is 2 spaces + tab
   	*              this is 3 spaces + tab
    	*          this is 4 spaces + tab
     	*          this is 5 spaces + tab
      	*          this is 6 spaces + tab
       	*          this is 7 spaces + tab
        	*      this is 8 spaces + tab
         	*      this is 9 spaces + tab
          	*      this is 10 spaces + tab
           	*      this is 11 spaces + tab
            	*  this is 12 spaces + tab */
    export class c4 {
    }
}

//// [commentsFormatting.js]
"use strict";
var m;
(function (m) {
    /** this is first line - aligned to class declaration
* this is 4 spaces left aligned
 * this is 3 spaces left aligned
  * this is 2 spaces left aligned
   * this is 1 spaces left aligned
    * this is at same level as first line
     * this is 1 spaces right aligned
      * this is 2 spaces right aligned
       * this is 3 spaces right aligned
        * this is 4 spaces right aligned
         * this is 5 spaces right aligned
          * this is 6 spaces right aligned
           * this is 7 spaces right aligned
            * this is 8 spaces right aligned */
    class c {
    }
    m.c = c;
    /** this is first line - 4 spaces right aligned to class but in js file should be aligned to class declaration
* this is 8 spaces left aligned
* this is 7 spaces left aligned
* this is 6 spaces left aligned
* this is 5 spaces left aligned
* this is 4 spaces left aligned
 * this is 3 spaces left aligned
  * this is 2 spaces left aligned
   * this is 1 spaces left aligned
    * this is at same level as first line
     * this is 1 spaces right aligned
      * this is 2 spaces right aligned
       * this is 3 spaces right aligned
        * this is 4 spaces right aligned
         * this is 5 spaces right aligned
          * this is 6 spaces right aligned
           * this is 7 spaces right aligned
            * this is 8 spaces right aligned */
    class c2 {
    }
    m.c2 = c2;
    /** this is comment with new lines in between

this is 4 spaces left aligned but above line is empty

 this is 3 spaces left aligned but above line is empty

  this is 2 spaces left aligned but above line is empty

    this is 1 spaces left aligned but above line is empty

    this is at same level as first line but above line is empty

     this is 1 spaces right aligned but above line is empty

      this is 2 spaces right aligned but above line is empty

       this is 3 spaces right aligned but above line is empty

        this is 4 spaces right aligned but above line is empty
    
    
    Above 2 lines are empty
    
    
    
    above 3 lines are empty*/
    class c3 {
    }
    m.c3 = c3;
    /** this is first line - aligned to class declaration
    *              this is 0 space + tab
    *              this is 1 space + tab
    *              this is 2 spaces + tab
    *              this is 3 spaces + tab
        *          this is 4 spaces + tab
        *          this is 5 spaces + tab
        *          this is 6 spaces + tab
        *          this is 7 spaces + tab
            *      this is 8 spaces + tab
            *      this is 9 spaces + tab
            *      this is 10 spaces + tab
            *      this is 11 spaces + tab
                *  this is 12 spaces + tab */
    class c4 {
    }
    m.c4 = c4;
})(m || (m = {}));


//// [commentsFormatting.d.ts]
namespace m {
    /** this is first line - aligned to class declaration
* this is 4 spaces left aligned
 * this is 3 spaces left aligned
  * this is 2 spaces left aligned
   * this is 1 spaces left aligned
    * this is at same level as first line
     * this is 1 spaces right aligned
      * this is 2 spaces right aligned
       * this is 3 spaces right aligned
        * this is 4 spaces right aligned
         * this is 5 spaces right aligned
          * this is 6 spaces right aligned
           * this is 7 spaces right aligned
            * this is 8 spaces right aligned */
    class c {
    }
    /** this is first line - 4 spaces right aligned to class but in js file should be aligned to class declaration
* this is 8 spaces left aligned
* this is 7 spaces left aligned
* this is 6 spaces left aligned
* this is 5 spaces left aligned
* this is 4 spaces left aligned
 * this is 3 spaces left aligned
  * this is 2 spaces left aligned
   * this is 1 spaces left aligned
    * this is at same level as first line
     * this is 1 spaces right aligned
      * this is 2 spaces right aligned
       * this is 3 spaces right aligned
        * this is 4 spaces right aligned
         * this is 5 spaces right aligned
          * this is 6 spaces right aligned
           * this is 7 spaces right aligned
            * this is 8 spaces right aligned */
    class c2 {
    }
    /** this is comment with new lines in between

this is 4 spaces left aligned but above line is empty

 this is 3 spaces left aligned but above line is empty

  this is 2 spaces left aligned but above line is empty

    this is 1 spaces left aligned but above line is empty

    this is at same level as first line but above line is empty

     this is 1 spaces right aligned but above line is empty

      this is 2 spaces right aligned but above line is empty

       this is 3 spaces right aligned but above line is empty

        this is 4 spaces right aligned but above line is empty
    
    
    Above 2 lines are empty
    
    
    
    above 3 lines are empty*/
    class c3 {
    }
    /** this is first line - aligned to class declaration
    *              this is 0 space + tab
    *              this is 1 space + tab
    *              this is 2 spaces + tab
    *              this is 3 spaces + tab
        *          this is 4 spaces + tab
        *          this is 5 spaces + tab
        *          this is 6 spaces + tab
        *          this is 7 spaces + tab
            *      this is 8 spaces + tab
            *      this is 9 spaces + tab
            *      this is 10 spaces + tab
            *      this is 11 spaces + tab
                *  this is 12 spaces + tab */
    class c4 {
    }
}


//// [DtsFileErrors]


commentsFormatting.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== commentsFormatting.d.ts (1 errors) ====
    namespace m {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /** this is first line - aligned to class declaration
    * this is 4 spaces left aligned
     * this is 3 spaces left aligned
      * this is 2 spaces left aligned
       * this is 1 spaces left aligned
        * this is at same level as first line
         * this is 1 spaces right aligned
          * this is 2 spaces right aligned
           * this is 3 spaces right aligned
            * this is 4 spaces right aligned
             * this is 5 spaces right aligned
              * this is 6 spaces right aligned
               * this is 7 spaces right aligned
                * this is 8 spaces right aligned */
        class c {
        }
        /** this is first line - 4 spaces right aligned to class but in js file should be aligned to class declaration
    * this is 8 spaces left aligned
    * this is 7 spaces left aligned
    * this is 6 spaces left aligned
    * this is 5 spaces left aligned
    * this is 4 spaces left aligned
     * this is 3 spaces left aligned
      * this is 2 spaces left aligned
       * this is 1 spaces left aligned
        * this is at same level as first line
         * this is 1 spaces right aligned
          * this is 2 spaces right aligned
           * this is 3 spaces right aligned
            * this is 4 spaces right aligned
             * this is 5 spaces right aligned
              * this is 6 spaces right aligned
               * this is 7 spaces right aligned
                * this is 8 spaces right aligned */
        class c2 {
        }
        /** this is comment with new lines in between
    
    this is 4 spaces left aligned but above line is empty
    
     this is 3 spaces left aligned but above line is empty
    
      this is 2 spaces left aligned but above line is empty
    
        this is 1 spaces left aligned but above line is empty
    
        this is at same level as first line but above line is empty
    
         this is 1 spaces right aligned but above line is empty
    
          this is 2 spaces right aligned but above line is empty
    
           this is 3 spaces right aligned but above line is empty
    
            this is 4 spaces right aligned but above line is empty
        
        
        Above 2 lines are empty
        
        
        
        above 3 lines are empty*/
        class c3 {
        }
        /** this is first line - aligned to class declaration
        *              this is 0 space + tab
        *              this is 1 space + tab
        *              this is 2 spaces + tab
        *              this is 3 spaces + tab
            *          this is 4 spaces + tab
            *          this is 5 spaces + tab
            *          this is 6 spaces + tab
            *          this is 7 spaces + tab
                *      this is 8 spaces + tab
                *      this is 9 spaces + tab
                *      this is 10 spaces + tab
                *      this is 11 spaces + tab
                    *  this is 12 spaces + tab */
        class c4 {
        }
    }
    