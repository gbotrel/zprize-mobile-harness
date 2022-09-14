package com.example.zprize;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

public class GnarkMSM {

    // note: on 32bit devices need to fallback to another method. not the case for the ZPrize.

    public String runMSMRandomMultipleVecs(String dir, String iters, String numElems, String numVecs) {
        String[] cmd = new String[]{gnarkPath, "-i", iters, "-n", numElems, "-v", numVecs, "-wd", dir};
        return runGnark(cmd);
    }

    public String runMSMFile(String dir, String iters) {
        String[] cmd = new String[]{gnarkPath, "-i", iters, "-t", "-wd", dir};
        return runGnark(cmd);
    }

    private String runGnark(String[] cmd) {
        String result = "";
        try {
            Process proc = Runtime.getRuntime().exec(cmd);

            BufferedReader stdInput = new BufferedReader(new
                    InputStreamReader(proc.getInputStream()));

            BufferedReader stdError = new BufferedReader(new
                    InputStreamReader(proc.getErrorStream()));

            // STDOUT
            String s = null;
            while ((s = stdInput.readLine()) != null) {
                System.out.println(s);
                result += s;
            }

            // STDERR
            while ((s = stdError.readLine()) != null) {
                System.out.println(s);
                result += s;
            }
            proc.destroy();
        } catch (IOException e) {
            result = e.toString();
            e.printStackTrace();
        }

        String finalResult = result;
        return finalResult;
    }


    private static String gnarkPath = "/data/data/com.example.zprize/lib/lib_gnark_.so";
}
